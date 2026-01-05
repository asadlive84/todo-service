package main

import (
	"context"
	"database/sql"

	// "database/sql"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	graphqlH "todo-service/internal/api/graphql/handler"
	"todo-service/internal/background/worker"

	// "todo-service/internal/background/eventconsumer"
	"todo-service/internal/background/eventconsumer"
	// pe "todo-service/internal/background/publish-event"
	"todo-service/internal/infrastructure/migration"

	// "todo-service/internal/infrastructure/repository"
	"todo-service/internal/infrastructure/storage"
	fileUseCase "todo-service/internal/usecase/file"
	todoUseCase "todo-service/internal/usecase/todo"

	"git.ice.global/packages/beeorm/v4"
	"git.ice.global/packages/hitrix"
	"git.ice.global/packages/hitrix/pkg/middleware"
	"git.ice.global/packages/hitrix/service"
	"git.ice.global/packages/hitrix/service/component/app"
	"git.ice.global/packages/hitrix/service/registry"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/vektah/gqlparser/v2/ast"

	beeORMentity "todo-service/internal/repository/beeorm"
	"todo-service/internal/repository/beeorm/entity"
	beeORMRepo "todo-service/internal/repository/beeorm/repository"
	fileRepo "todo-service/internal/repository/file"
	redisSearch "todo-service/internal/repository/search/redis"

	_ "github.com/go-sql-driver/mysql"
)

func cmd() {
	// Initialize hitrix server
	s, deferFunc := hitrix.New(
		"todo-app", "your secret",
	).RegisterDIGlobalService(
		registry.ServiceProviderConfigDirectory("./config"), // without docker
		registry.ServiceProviderUUID(),
		registry.ServiceProviderErrorLogger(),
		// registry.ServiceProviderConfigDirectory("/app/config"), for docker
		registry.ServiceProviderOrmRegistry(beeORMentity.Init),
		registry.ServiceProviderOrmEngine(),
		registry.ServiceProviderJWT(),

		// customService.ServiceProviderRedisSearch(),
	).RegisterDIRequestService(
		registry.ServiceProviderOrmEngineForContext(true),
	).RegisterRedisPools(&app.RedisPools{Persistent: "persistent"}).
		Build()

	defer deferFunc()

	configService := service.DI().Config()

	DRIVER_NAME := configService.DefString("DATABASE.MYSQL.DRIVER_NAME", "mysql")
	USER_NAME := configService.DefString("DATABASE.MYSQL.USER_NAME", "root")
	PASSWORD := configService.DefString("DATABASE.MYSQL.PASSWORD", "password")
	HOST := configService.DefString("DATABASE.MYSQL.HOST", "localhost")
	PORT := configService.DefString("DATABASE.MYSQL.PORT", "3306")
	DATABASE_NAME := configService.DefString("DATABASE.MYSQL.DATABASE_NAME", "hitrix_test")

	APP_PORT := configService.DefString("server.port", "8080")

	// REDIS_ADDR := configService.DefString("REDIS.REDIS_ADDR", "localhost:6379")
	// REDIS_STREAM := configService.DefString("REDIS.REDIS_STREAM", "todos:events")

	s3Endpoint := configService.DefString("S3.S3_ENDPOINT", "http://localstack:4566")
	s3Bucket := configService.DefString("S3.S3_BUCKET", "todo-bucket")

	migrationDSN := fmt.Sprintf(
		"%s://%s:%s@tcp(%s:%s)/%s",
		DRIVER_NAME,
		USER_NAME,
		PASSWORD,
		HOST,
		PORT,
		DATABASE_NAME,
	)

	migrationsPath := "./migrations"
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		migrationsPath = "/app/migrations"
	}

	log.Info().Msgf("migrationsPath: %+v", migrationsPath)

	if err := migration.RunMigrations(migrationDSN, migrationsPath); err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}
	log.Info().Msg("Migrations completed successfully")

	ormengine := service.DI().OrmEngine()

	// schema := ormengine.GetRegistry().GetTableSchemaForEntity(&entity.TodoEntity{})

	// search := ormengine.GetRedisSearch("todo_cache")
	// search.Query("FT.DROPINDEX", "entity.TodoEntity")

	// schema.UpdateSchema(ormengine)

	// schema.ReindexRedisSearchIndex(ormengine)

	fmt.Println("RedisSearch index is ready!")

	alters := ormengine.GetAlters()
	for _, alter := range alters {
		log.Printf("Executing Alter on [%s]: %s", alter.Pool, alter.SQL)
		alter.Exec()
	}

	dsn := fmt.Sprintf(

		"%s:%s@tcp(%s:%s)/%s?parseTime=true",

		DRIVER_NAME,

		PASSWORD,

		HOST,

		PORT,

		DATABASE_NAME,
	)

	db, err := sql.Open("mysql", dsn)

	if err != nil {

		log.Fatal().Err(err).Msg("Failed to open database")

	}

	defer db.Close()

	db.SetMaxOpenConns(25)

	db.SetMaxIdleConns(5)

	db.SetConnMaxLifetime(5 * time.Minute)

	processor := worker.NewOutboxProcessor(ormengine)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go processor.Start(ctx)

	// StartBackgroundWorker(ormengine)
	fmt.Println("==========InitSearchIndex=============")
	// InitSearchIndex2(ormengine)
	InitSearchIndex(ormengine)
	RunSearchConsumer(ormengine)

	// Initialize repositories
	todoRepo := beeORMRepo.NewOrmEngine(ormengine)
	fileRepo := fileRepo.NewFileRepository(db) // MUST CHANGE

	log.Info().Msgf("S3_BUCKET=%s, S3_ENDPOINT=%s", s3Bucket, s3Endpoint)

	s3Repo, err := storage.NewS3Repository(context.Background(), s3Bucket, s3Endpoint)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize S3")
	}

	// redis.InitRedis(REDIS_ADDR)

	newRedisSearch := redisSearch.NewRedisSearchService(ormengine)

	// Initialize use cases
	todoUC := todoUseCase.NewTodoUseCase(todoRepo, s3Repo, newRedisSearch, s3Bucket)
	fileUC := fileUseCase.NewFileUseCase(fileRepo, todoRepo, s3Repo, s3Bucket)

	// Convert port to uint for hitrix
	portNum, err := strconv.ParseUint(APP_PORT, 10, 32)
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid APP_PORT value")
	}

	// Setup GraphQL schema
	executableSchema := graphqlH.NewExecutableSchema(graphqlH.Config{
		Resolvers: &graphqlH.Resolver{
			TodoUseCase: todoUC,
			FileUseCase: fileUC,
		},
	})

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info().Msg("Starting background processes...")
		s.RunBackgroundProcess(func(b *hitrix.BackgroundProcessor) {
			// b.RunScript(&pe.PublishEvent{})
			b.RunScript(&eventconsumer.EventConsumer{})
		})
	}()

	// Start server in goroutine
	go func() {
		log.Info().Msgf("GraphQL server starting on :%s", APP_PORT)
		log.Info().Msgf("GraphQL Playground available at http://localhost:%s/playground", APP_PORT)

		s.RunServer(uint(portNum), executableSchema, func(ginEngine *gin.Engine) {
			middleware.Cors(ginEngine)

			// Print all registered routes
			for _, route := range ginEngine.Routes() {
				log.Info().Msgf("Route: %s %s", route.Method, route.Path)
			}

			ginEngine.GET("/playground", gin.WrapH(playground.Handler("GraphQL Playground", "/query")))

			ginEngine.GET("/health", func(c *gin.Context) {
				c.JSON(200, gin.H{"status": "healthy"})
			})
		}, gqlSetup)

	}()

	// Wait for shutdown signal
	<-quit
	log.Info().Msg("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Add any cleanup logic here if needed
	select {
	case <-ctx.Done():
		log.Info().Msg("Shutdown timeout exceeded")
	default:
		log.Info().Msg("Server shutdown successfully")
	}
}

func gqlSetup(srv *handler.Server) {
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

}

// func InitSearchIndex(engine *beeorm.Engine) {
// 	schema := engine.GetRegistry().GetTableSchemaForEntity(&entity.TodoEntity{})

// 	schema.ReindexRedisSearchIndex(engine)

// 	log.Print("RedisSearch Index has been re-created automatically!")

// 	engine.GetRegistry().GetTableSchemaForEntity(&entity.TodoEntity{}).ReindexRedisSearchIndex(engine)

// 	log.Print("RedisSearch index is ready!")
// }

func InitSearchIndex(engine *beeorm.Engine) {
	schema := engine.GetRegistry().GetTableSchemaForEntity(&entity.TodoEntity{})

	has, alters := schema.GetSchemaChanges(engine)
	if has {
		for _, alter := range alters {
			log.Printf("🛠️ Applying Index Alter: %s", alter.SQL)
			alter.Exec()
		}
	}

	schema.ReindexRedisSearchIndex(engine)

	log.Print("✅ RedisSearch index is synced and ready!")
}
func RunSearchConsumer(engine *beeorm.Engine) {
	consumer := beeorm.NewBackgroundConsumer(engine)
	go func() {
		for {
			consumer.Digest(context.Background())
			time.Sleep(time.Second)
		}
	}()
}
