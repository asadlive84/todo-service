.PHONY: run test benchmark generate-mocks

#docker compose down --remove-orphans -v && docker compose up --build -d
run:
	docker compose up --build

test:
	go test ./test/unit/... -v

benchmark:
	go test ./test/bench/... -bench=. -benchmem

generate-mocks:
	mockgen -source=internal/domain/interface/interface.go -destination=test/mock/genrate-mocks.go -package=mock todo-service/internal/domain/interface