# LOGISTIC APP

<img src="https://img.shields.io/badge/Go-1.22.0-00ADD8?style=for-the-badge&logo=go" alt="go version" />

![Build Status](https://github.com/mahsa-fathi/logistic-app/actions/workflows/ci-build.yaml/badge.svg)

A service for connecting customers to providers and updating order status based on an API provided by the provider.

This project uses hexagonal architecture.

## ⚡️ Quick start

An easy way to start the project is to use docker. Server is then available on port 8080:

```shell
docker-compose up -d
```

Another way to start the project is: 

1. start a postgres database
2. install mod
    ```shell
    go mod download
    ```
3. provide the following environment variables for database connection
    ```shell
    DB_ADDRESS=localhost
   DB_NAME=postgres
   DB_PASSWORD=pg_pass
   DB_PORT=5432
   DB_USER=pg
    ```
4. run the following command to start the server
   ```shell
   go build -o server ./cmd/http
   ./server
    ```
5. run the following command to start the cron jobs
    ```shell
    go build -o cron ./cmd/cron
   ./cron
    ```

## 🗄 Template structure

### ./internal/adapters

**Folder with all adapters**. This directory is used for both driven and driver adapters, including repositories and http servers. 

- `./internal/adapters/cron` folder for running a scheduler
- `./internal/adapters/db` folder for connection and queries to database
- `./internal/adapters/http` folder for running a http server


### ./internal/app

**Folder with application logic only**. This directory doesn't care about _what database driver you're using_ or _which caching solution your choose_ or any third-party things.

- `./internal/app/domain` folder for describing models
- `./internal/app/ports` folder for describing ports for driven and driver interfaces
- `./internal/app/service` folder for functional usage

### ./internal/common

**Folder for common features**. This directory contains everything that may be used across the project.

- `./internal/common/configs` folder for importing every configuration like environment variable
- `./internal/common/errors` folder for errors across the project

## ⚙️ Configuration

```ini
# .env

# Server settings:
SERVER_URL=localhost:8080
SERVER_READ_TIMEOUT=60

# JWT settings:
SECRET_KEY=secret
TOKEN_EXPIRATION=24  #in hours

# Logging
LOG_ERROR="yes"

# Database
DB_ADDRESS=localhost
DB_NAME=postgres
DB_PASSWORD=pg_pass
DB_PORT=5432
DB_USER=pg

# Test Database
DB_TEST_ADDRESS=localhost
DB_TEST_NAME=postgres
DB_TEST_PASSWORD=pg_pass
DB_TEST_PORT=5432
DB_TEST_USER=pg

# Periodic Tasks
ORDER_UPDATE_PERIOD=86400  #in seconds
PERIODIC_TASK_MAX_CONCURRENCY=10  #concurrency of running goroutines for updating order status
```

## 📦 Data Model

### Customers

Represents a customer involved in an order (either as sender or receiver).
Sender can initialize an order.

| Field       | Type      | Description                   |
|-------------|-----------|-------------------------------|
| ID          | uint      | Primary key (auto-increment). |
| PhoneNumber | string    | Unique                        |
| Name        | string    | Optional                      |
| Address     | string    | not null                      |
| PostalCode  | string    | not null                      |
| CreatedAt   | Timestamp |                               |
| UpdatedAt   | Timestamp |                               |

### Providers

Represents the service provider responsible for the delivery.

| Field     | Type      | Description                   |
|-----------|-----------|-------------------------------|
| ID        | uint      | Primary key (auto-increment). |
| Name      | string    | unique                        |
| Url       | string    | not null                      |
| CreatedAt | Timestamp |                               |
| UpdatedAt | Timestamp |                               |

### Orders

This table keeps the data for any order that a user is registered. 
It also keeps the data for when the package was picked up and when it was delivered.
This data can be used for when we want a report of providers delivery time.
A NotifiedReceiver column is kept to avoid multiple messages to the receiver. 
Default value for status is set to PROVIDER_SEEN in order to facilitate the program.

Every id in this table is indexed. Besides, created_at is another index for when we want to get reports of the this table.
If a very bulky table is anticipated we can also use partitioning created_at column on this table for better performance, which is not considered in this project.

A partial index is used for this table with following format. This partial index is used for identifying ongoing orders.

```sql
CREATE INDEX CONCURRENTLY idx_ongoing_status
   ON orders(status)
   WHERE status IN ('IN_PROGRESS', 'PROVIDER_SEEN', 'PICKED_UP');
```

| Field            | Type        | Description                    |
|------------------|-------------|--------------------------------|
| ID               | uint        | Primary key (auto-increment).  |
| Provider         | Provider    | Foreign Key to providers table |
| Sender           | Customer    | Foreign Key to customers table |
| Receiver         | Customer    | Foreign Key to customers table |
| Product          | string      | Optional product description.  |
| Status           | varchar(15) | Default: 'PROVIDER_SEEN'       |
| PickedUpDate     | Date        |                                |
| DeliveryDate     | Date        |                                |
| NotifiedReceiver | bool        | default is false.              |
| CreatedAt        | Timestamp   |                                |
| UpdatedAt        | Timestamp   |                                |

### PeriodicTasks

This table keep tracks of the cron jobs ran through the program. A way to see errors and if they were successful.

| Field            | Type      | Description                          |
|------------------|-----------|--------------------------------------|
| ID               | uint      | Primary key (auto-increment).        |
| JobName          | string    | unique                               |
| IntervalInMinute | uint      | not null                             |
| LastRunTime      | Timestamp |                                      |
| Failed           | bool      |                                      |
| Error            | string    |                                      |
| CreatedAt        | Timestamp | Timestamp when the task was created. |

## 🔍 API Endpoints

### GET /api/providers/

Returns a list of all registered providers.

```shell
curl -X GET http://localhost:8080/api/providers/
```

Example response:

```json
[
    {
        "id": 1,
        "name": "test-provider-1",
        "url": "https://staging.podro.com/api/mock/status",
        "created_at": "2025-04-23T19:20:21.979013+03:30",
        "updated_at": "2025-04-23T19:20:21.979013+03:30"
    },
    {
        "id": 2,
        "name": "test-provider-2",
        "url": "https://staging.podro.com/api/mock/status",
        "created_at": "2025-04-24T14:26:19.58349+03:30",
        "updated_at": "2025-04-24T14:26:19.58349+03:30"
    }
]
```

### GET /api/providers/report/

Returns the average delivery time (in days) for each provider over the past 7 days in a descending order.

```shell
curl -X GET http://localhost:8080/api/providers/report/
```

Example response:

```json
[
    {
        "provider_id": 1,
        "mean_delivery_time_in_days": 5
    },
    {
        "provider_id": 2,
        "mean_delivery_time_in_days": 3
    }
]
```

### POST /api/provider/

Creates a new provider.

```shell
curl -X POST http://localhost:8080/api/provider/ \
  -H "Content-Type: application/json" \
  -d '{"name": "test-provider-3", "url": "https://staging.podro.com/api/mock/status"}'
```

Example response:
```json
{
    "id": 5,
    "name": "test-provider-3",
    "url": "https://staging.podro.com/api/mock/status",
    "created_at": "2025-04-25T02:41:54.1687205+03:30",
    "updated_at": "2025-04-25T02:41:54.1687205+03:30"
}
```

### POST /api/customer/

Registers a new customer.

```shell
curl -X POST 'http://localhost:8080/api/customer/' \
-H 'Content-Type: application/json' \
-d '{
    "name": "mahsa",
    "phone_number": "09378",
    "address": "somewhere",
    "postal_code": "6372687"
}'
```

Example response

```json
{
    "id": 8,
    "phone_number": "09378",
    "name": "mahsa",
    "address": "somewhere",
    "postal_code": "6372687",
    "created_at": "2025-04-25T02:43:59.9970862+03:30",
    "updated_at": "2025-04-25T02:43:59.9970862+03:30"
}
```

### POST /api/customer/token/

Retrieves a token for an existing customer. Must be used to get or create orders.

```shell
curl -X POST http://localhost:8080/api/customer/token/ \
  -H "Content-Type: application/json" \
  -d '{"id": 8}'
```

Example response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5..."
}
```

### POST /api/order/

Creates a new order. Requires authentication.

```shell
curl -X POST http://localhost:8080/api/order/ \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "provider_id": 1,
    "receiver_id": 2,
    "product": "Books"
  }'
```

Example response:

```json
{
    "id": 3,
    "provider_id": 2,
    "sender_id": 5,
    "receiver_id": 8,
    "product": "book",
    "status": "PENDING",
    "picked_up_date": null,
    "delivery_date": null,
    "notified_receiver": false,
    "created_at": "2025-04-25T02:47:29.4592826+03:30",
    "updated_at": "2025-04-25T02:47:29.4592826+03:30"
}
```

### GET /api/order/{order_id}/

Fetches an order by ID. Requires authentication. Must be sender or receiver.

```shell
curl -X GET http://localhost:8080/api/order/3/ \
  -H "Authorization: Bearer <TOKEN>"
```

Example response:

```json
{
    "id": 3,
    "provider_id": 2,
    "provider": {
        "id": 2,
        "name": "test-provider-2",
        "url": "https://staging.podro.com/api/mock/status",
        "created_at": "2025-04-24T14:26:19.58349+03:30",
        "updated_at": "2025-04-24T14:26:19.58349+03:30"
    },
    "sender_id": 5,
    "sender": {
        "id": 5,
        "phone_number": "09357",
        "name": "mahsa",
        "address": "somewhere",
        "postal_code": "6372687",
        "created_at": "2025-04-24T14:30:33.395016+03:30",
        "updated_at": "2025-04-24T14:30:33.395016+03:30"
    },
    "receiver_id": 8,
    "receiver": {
        "id": 8,
        "phone_number": "09378",
        "name": "mahsa",
        "address": "somewhere",
        "postal_code": "6372687",
        "created_at": "2025-04-25T02:43:59.997086+03:30",
        "updated_at": "2025-04-25T02:43:59.997086+03:30"
    },
    "product": "book",
    "status": "DELIVERED",
    "picked_up_date": "2025-04-22T00:00:00Z",
    "delivery_date": "2025-04-23T00:00:00Z",
    "notified_receiver": true,
    "created_at": "2025-04-25T02:47:29.459282+03:30",
    "updated_at": "2025-04-25T02:52:43.376242+03:30"
}
```

## ⏱️ Cron Jobs

There is only one cron job in this project that runs each 24 hours to update the status of each order. 
It only runs on ongoing orders. The code is located in `internal/app/service/order_task.go`.

To test this part, `ORDER_UPDATE_PERIOD` environment variable can be used to reduce the interval of this periodic task (It is set in seconds).
A random choice is used to update the status of products based on the mocked url given in the project description.
If there are any failures, code is retried 3 times and then logs the error on the periodic_tasks table. 

In this periodic task the receiver is notified if the status of the url indicates `PICKED_UP`. 
First a query is called on the orders table, to see if the receiver was already notified or not, then a message is sent to them.
