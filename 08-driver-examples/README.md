## Run PostgreSQL in Docker
```shell
docker run -d --name data-persistence -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -e POSTGRES_DB=data-persistence postgres:15
```

## Migrate
Replace `project_path` according your needs.

### Migrate up
```shell
migrate -path '<project_path>/data-persistence/migrations' -database 'postgres://root:root@localhost:5432/data-persistence?sslmode=disable' up
```

### Migrate down
```shell
migrate -path '<project_path>/data-persistence/migrations' -database 'postgres://root:root@localhost:5432/data-persistence?sslmode=disable' down
```
