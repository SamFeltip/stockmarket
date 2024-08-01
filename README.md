Stockmarket
===========

Stockmarket is a stock market simulator built with golang, HTMX, and PostgreSQL.

it is self hosted using nginx, docker-compose, and docker.

## Features

- Stocks
- Insights
- Stocks can be bought or sold
- Games can be created and played
- After each round in a game (once every player has made 3 moves) the stock prices is updated according to player insights

## Repo Structure

- controllers - contains the controllers for the routes
- database - contains the database schema and migrations
- frontend - contains the frontend code
- models - contains the models for the database
- pkg - contains the packages used by the project

## Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Running the project

`cd backend/`
`go run .`

postgres should be running on port 5433 (this is specified in a .env file)

### .env layout

```
JWT_SECRET=secret
PORT=4040
DB_HOST=localhost
DB_USER=me
DB_PASSWORD=password

DB_NAME=postgres
DB_PORT=5433

ENVIRONMENT=development
```

in deployment:

```
ENVIRONMENT=production
DB_HOST=db
DB_USER=postgres
```

Images are kept out of git for the skae of push/pull performance. locally these are kept at `~/Documents/Code/Decent Projects/stockmarket/backend/static/imgs` on my laptop.

### seeding

`setup.sql` contains the seeding for the database. this includes insight templates and stocks.
