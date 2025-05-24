1. Install Go
Download and install Go from the official website:
👉 https://go.dev/dl/

Ensure it's properly installed by running:
go version


2. Install Docker and Docker Compose (macos/windows etc)
Download Docker Desktop from:
👉 https://www.docker.com/products/docker-desktop/

Make sure both Docker and Docker Compose are working:
docker --version
docker-compose --version


3. Install Goose (Database Migration Tool)
Goose is a tool used for managing database migrations in Go. Install it using:
go install github.com/pressly/goose/v3/cmd/goose@latest
For more info, visit the Goose GitHub repo:
👉 https://github.com/pressly/goose


4. Start PostgreSQL with Docker
Once all the above tools are installed, spin up the PostgreSQL container:
"docker-compose up -d"
This will start a single PostgreSQL instance using the configuration defined in your docker-compose.yml file.


5. Apply Database Migrations
Run the set_env.sh script to configure the necessary environment variables for Goose, then, apply the migrations using:
"export GOOSE_DRIVER=postgres"
"export GOOSE_DBSTRING=postgres://alle:password@localhost:5432/alle_prod?sslmode=disable"
"goose up" (This will create the required task table in your local database.)


6. Run the Server
Start the Go application from the root directory:
"go run src/main.go"  -> it will start listening at 8080

7. Start firing Api's

Create Task Api-
Request:
curl --location 'localhost:8080/v1/task' \
--header 'Content-Type: application/json' \
--data '{
    "name": "task5",
    "status": "COMPLETED"
}'

Response:
{
    "data": {
        "Id": "f9004468-1e84-4512-9240-11385395d132"
    },
    "meta": null,
    "status_code": 200
}

Update Task Api-

Request:
curl --location --request PATCH 'localhost:8080/v1/task/f5134c46-409b-418c-bbb7-1008fee4af3b' \
--header 'Content-Type: application/json' \
--data '{
    "name": "task1",
    "status": "PENDING"
}'

Response:
{
    "data": "f5134c46-409b-418c-bbb7-1008fee4af3b",
    "meta": null,
    "status_code": 200
}


Get Tasks Api-

Request:
curl --location 'localhost:8080/v1/tasks?page=3&per_page=2&sort=id'

Response:
{
    "data": [
        {
            "id": "f9004468-1e84-4512-9240-11385395d132",
            "name": "task5",
            "status": "COMPLETED",
            "created_at": "2025-05-24T17:21:21Z",
            "modified_at": "2025-05-24T17:21:21Z"
        }
    ],
    "meta": {
        "total": 5,
        "page": 3,
        "per_page": 2
    },
    "status_code": 200
}


Delete Task-

Request:
curl --location --request DELETE 'localhost:8080/v1/task/f5134c46-409b-418c-bbb7-1008fee4af3b'

Response:
{
    "data": "f5134c46-409b-418c-bbb7-1008fee4af3b",
    "meta": null,
    "status_code": 200
}