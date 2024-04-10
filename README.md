
## Game Blitz API Documentation

### Overview

The Game Blitz API is designed to manage basic gaming features such as Statistics, Quests, and Leaderboards. It is implemented using Golang version 1.22 and provides a robust set of endpoints for creating, retrieving, updating, and deleting game-related data.

### Features

- **Leaderboards**: Create, retrieve, update, and delete leaderboards.
- **Quests**: Manage quests and their associated tasks.
- **Statistics**: Handle player statistics and track progress.
- **Player Progression**: Track and update player progress in quests and statistics.

### Prerequisites

- Golang version 1.21

### Installation

1. Clone the repository:
```bash
git clone github.com/gabapcia/gameblitz
```

2. Navigate to the project directory:
```bash
cd game-blitz
```

### Configuration

Before running the application, set the required environment variables:

| Variable                         | Description                                      | Type    | Required | Example                                                                   |
|----------------------------------|--------------------------------------------------|---------|----------|---------------------------------------------------------------------------|
| `PORT`                           | API Port to listen to                            | Integer | Yes      | `8080`                                                                    |
| `KEYCLOACK_CERTS_URI`            | Keycloack certs URI                              | String  | Yes      | `http://localhost:3000/realms/gameblitz/protocol/openid-connect/certs`    |
| `POSTGRESQL_DSN`                 | PostgreSQL connection string                     | String  | Yes      | `postgres://gameblitz:gameblitz@localhost:5432/gameblitz?sslmode=disable` |
| `MONGO_URI`                      | MongoDB connection string                        | String  | Yes      | `mongodb://localhost:27017/?retryWrites=true&w=majority`                  |
| `MONGO_DB`                       | MongoDB database name                            | String  | Yes      | `gameblitz`                                                               |
| `REDIS_ADDR`                     | Redis address                                    | String  | Yes      | `localhost:6379`                                                          |
| `REDIS_USERNAME`                 | Redis username                                   | String  | No       | `gameblitz`                                                               |
| `REDIS_PASSWORD`                 | Redis password                                   | String  | No       | `gameblitz`                                                               |
| `REDIS_DB`                       | Redis database                                   | Integer | No       | `0`                                                                       |
| `MEMCACHED_CONN_STR`             | Memcached connection string                      | String  | Yes      | `localhost:11211`                                                         |
| `MEMCACHED_EXPIRATION`           | Cache expiration in seconds for the GET endpoint | Integer | No       | `60`                                                                      |
| `MEMCACHED_MIDDLEWARE_EXPIRATION`| Cache expiration in seconds for the Middlewares  | Integer | No       | `60`                                                                      |
| `RABBITMQ_URI`                   | RabbitMQ connection string                       | String  | Yes      | `amqp://gameblitz:gameblitz@localhost:5672/gameblitz`                     |


### Running the Application

Build and start the application:

```bash
go build -o game-blitz cmd/api/main.go
./game-blitz
```

To explore the API, navigate to the `/docs` endpoint where the Swagger documentation is available.

### Running Tests

To execute the unit tests, run the following command:

```bash
make test
```

### License

This project is licensed under the MIT License.
