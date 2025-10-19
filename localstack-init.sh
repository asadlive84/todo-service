#!/bin/bash

set -e

# Check if S3_BUCKET is set
if [ -z "$S3_BUCKET" ]; then
  echo "Error: S3_BUCKET environment variable is not set"
  exit 1
fi

echo "Waiting for LocalStack S3 service to be ready..."
for i in {1..30}; do
  if awslocal s3 ls &>/dev/null; then
    echo "LocalStack S3 is ready!"
    break
  fi
  echo "Attempt $i: Waiting for S3..."
  sleep 2
done

# Create S3 bucket
BUCKET_NAME=$S3_BUCKET
echo "Creating S3 bucket: $BUCKET_NAME..."
awslocal s3 mb s3://$BUCKET_NAME 2>/dev/null || echo "Bucket $BUCKET_NAME already exists or creation skipped"

# Verify bucket creation
echo "Listing S3 buckets..."
awslocal s3 ls

echo "LocalStack initialization completed successfully!"