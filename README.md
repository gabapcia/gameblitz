
## Game Blitz API Documentation

### Overview

The Game Blitz API is designed to manage basic gaming features such as Statistics, Quests, and Leaderboards. It is implemented using Golang version 1.21 and provides a robust set of endpoints for creating, retrieving, updating, and deleting game-related data.

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
git clone github.com/gabarcia/game-blitz
```

2. Navigate to the project directory:
```bash
cd game-blitz
```

### Configuration

Before running the application, set the required environment variables:

```bash
export PORT=[port] # API Port to listen to, exemple: 8080
export POSTGRESQL_DSN=[postgresql_dsn] # PostgreSQL connection string, exemple: postgres://metagaming:metagaming@localhost:5432/metagaming?sslmode=disable
export MONGO_URI=[mongo_uri] # MongoDB connection string, exemple: mongodb://localhost:27017/?retryWrites=true&w=majority
export MONGO_DB=[mongo_db] # MongoDB database name, exemple: metagaming
export REDIS_ADDR=[redis_addr] # Redis address, exemple: localhost:6379
export REDIS_USERNAME=[redis_username] # Redis username, default: empty
export REDIS_PASSWORD=[redis_password] # Redis password, default: empty
export REDIS_DB=[redis_db] # Redis database, default: 0
export MEMCACHED_CONN_STR=[memcached_conn_str] # Memcached connection string
export MEMCACHED_EXPIRATION=[memcached_expiration] # Cache expiration in seconds for the GET endpoint, default: 60
export MEMCACHED_MIDDLEWARE_EXPIRATION=[memcached_middleware_expiration] # Cache expiration in seconds for the Middlewares, default: 60
export RABBITMQ_URI=[rabbitmq_uri] # RabbitMQ connection string, exemple: amqp://metagaming:metagaming@localhost:5672/metagaming
```

### Running the Application

Build and start the application:

```bash
go build -o game-blitz cmd/main.go
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
