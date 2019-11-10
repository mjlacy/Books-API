@echo off

IF "%1"=="func" (
   go test ./... -coverprofile coverage.out
   go tool cover -func=coverage.out
   del coverage.out
) ELSE IF "%1"=="html" (
   go test ./... -coverprofile coverage.out
   go tool cover -html=coverage.out
   del coverage.out
) ELSE (
   go test ./...
)
