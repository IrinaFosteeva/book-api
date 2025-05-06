# Book API

A simple REST API for managing books, built with Go and MongoDB.

## Setup

1. **Install Docker**: Ensure Docker is installed and running.
2. **Run MongoDB container**:
   ```bash
   docker run --rm -d -p 27017:27017 --name mongodb mongo:latest
   ```
3. **Clone the repository**:
   ```bash
   git clone https://github.com/IrinaFosteeva/book-api.git
   cd book-api
   ```
4. **Run the API**:
   ```bash
   docker run --rm -it -v $(pwd):/app -w /app -p 8080:8080 golang:1.21 bash
   ```
   Inside the container, run:
   ```bash
   go mod init book-api
   go get github.com/gorilla/mux
   go get go.mongodb.org/mongo-driver/mongo

   go get github.com/swaggo/swag@v1.16.4
   go get github.com/swaggo/http-swagger
   go install github.com/swaggo/swag/cmd/swag@v1.16.4
   swag init --dir ./cmd/api,./internal --output ./docs

   go run cmd/api/main.go
   ```

## Endpoints

- **POST /books**
  - Creates a new book.
  - Request body: `{"title":"string","author":"string"}`
  - Response: `201 Created`, book with generated `id`
  - Example:
    ```bash
    curl -X POST http://localhost:8080/books -d '{"title":"Book Title","author":"Author Name"}'
    ```
  - Example response:
    ```json
    {"id":"680fa8af4b925b6dd5c531c5","title":"Book Title","author":"Author Name"}
    ```

- **GET /books**
  - Returns all books.
  - Response: `200 OK`, array of books
  - Example:
    ```bash
    curl http://localhost:8080/books
    ```
  - Example response:
    ```json
    [
      {"id":"680fa8af4b925b6dd5c531c5","title":"Book Title","author":"Author Name"},
      {"id":"680fa8af4b925b6dd5c531c6","title":"Another Book","author":"Another Author"}
    ]
    ```

- **GET /books?id={id}**
  - Returns a book by ID.
  - Response: `200 OK`, book; `404 Not Found` if not found
  - Example:
    ```bash
    curl http://localhost:8080/books?id=680fa8af4b925b6dd5c531c5
    ```
  - Example response:
    ```json
    {"id":"680fa8af4b925b6dd5c531c5","title":"Book Title","author":"Author Name"}
    ```

- **DELETE /books?id={id}**
  - Deletes a book by ID.
  - Response: `204 No Content`; `404 Not Found` if not found
  - Example:
    ```bash
    curl -X DELETE http://localhost:8080/books?id=680fa8af4b925b6dd5c531c5
    ```

## Testing

Run unit tests:
```bash
docker run --rm -it -v $(pwd):/app -w /app golang:1.21 bash
go test ./tests
```

## Dependencies

- Go 1.21
- MongoDB
- github.com/gorilla/mux
- go.mongodb.org/mongo-driver/mongo