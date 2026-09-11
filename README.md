# Contact form

## Development with Compose

Start PostgreSQL, run the migrations and seed data, and launch the application:

```sh
docker compose up --build
```

The API is available at `http://localhost:8000`. You can change the published
ports with `APP_PORT` and `POSTGRES_PORT`:

```sh
APP_PORT=8080 POSTGRES_PORT=5433 docker compose up --build
```

To stop the environment while preserving the database:

```sh
docker compose down
```

## Dev Container

Open the repository using VS Code's **Dev Containers: Reopen in Container**
command. The container uses Go 1.27, downloads the dependencies, and starts
PostgreSQL automatically. Inside the container, start the service with:

```sh
go run . serve
```
