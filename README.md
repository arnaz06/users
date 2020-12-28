# Users

Users is service user management.

## API Documentation

Open the docs [docs/openapi.yaml](docs/openapi.yaml) with [Swagger Editor](http://editor.swagger.io/).

### Testing

There are two kind of tests.

#### Unit Test

```bash
make unittest
```

#### Integration Tests

Before running integration tests, make sure you have the right state `make migrate-up`

```bash
make mysql-up
make test
```

### Running

- Spin up the mysql database & run the migration Script.

```bash
make mysql-up
make migrate-up
```

- Run the API:

```bash
make docker
make run
make stop
```

Make sure to set the `.env` file (see: `.env.example`).
