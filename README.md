# Book API
REST API for managing books using Go and MongoDB.

## Setup
1. Run MongoDB container:
   ```bash
   docker run --rm -d -p 27017:27017 --name mongodb mongo:latest

2. Run Go container:
```bash
 docker run --rm -it -v ~/projects/book-api:/app -w /app -p 8080:8080 golang:1.21
 go mod init book-api
 go get github.com/gorilla/mux
 go get go.mongodb.org/mongo-driver/mongo
 go run cmd/api/main.go