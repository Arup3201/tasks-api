# Task API

## Set up

For setup, there are few software you need to run to make sure are running to run the API without any error. They are PostgreSQL and Keycloak. PostgreSQL is used as the database of the API and Keycloak is used for authentication-authorization.

In this doc, I am using docker to run all the essential applications. If you have them on you local machine then after running them the steps won't change.

### PostgreSQL with Docker

Create networks that can be used for the communication:

```shell
docker create network postgres-net
docker create network keycloak-net
```

`postgres-net` is the network to connect with the database. `keycloak-net` is the network to connect with keycloak server.

Following docker run will run a PostgreSQL server with the database `tests` that you will use for the API.

> You can change the database name if you want but make sure to change to other places where it is used as well.

```sh
docker run -d \
--network postgres-net \
--network-alias=postgres \
--volume psql-data:/var/lib/postgresql \
-e POSTGRES_PASSWORD=1234 \
-e POSTGRES_USERNAME=postgres \
-e POSTGRES_DB=tests \
postgres:18-alpine
```

### Keycloak with Docker

Next you need to start the keycloak service and attach it to PostgreSQL. You can do that with the following docker run command:

```sh
docker run \
-p 127.0.0.1:8080:8080 \
--network postgres-net \
--network keycloak-net --network-alias=keycloak \
-e KC_BOOTSTRAP_ADMIN_USERNAME=admin \
-e KC_BOOTSTRAP_ADMIN_PASSWORD=admin \
-e KC_DB=postgres \
-e KC_DB_URL=jdbc:postgresql://postgres/tests \
-e KC_DB_USERNAME=postgres \
-e KC_DB_PASSWORD=1234 \
quay.io/keycloak/keycloak:26.4.0 start-dev
```

You can open the keycloak at `127.0.0.1:8080`. Login using the username `admin` and password `admin`. Then proceed to create a new realm `tasks` and a new private client in that realm `api`.

If you are confused on how to do it, then you can refer this [gist](https://gist.github.com/upender-devulapally/22033a8f530acbe95696e3003de61eb3).

### API with Docker

Next, build the image of the API using docker:

```sh
docker build . -t tasks-api
```

And then, run the API by providing the environment variables inside a `.env` file like the following:

```sh
docker run -p 127.0.0.1:8086:8086 --network keycloak-net --network postgres-net --env-file .env tasks-api /tasks-api
```

The `.env` file should contain the following values:

```sh
DBHOST=postgres
DBUSER=postgres
DBPORT=5432
DBPASS=1234
DBNAME=tests
KEYCLOAK_SERVER_URL=http://keycloak:8080
KEYCLOAK_REALM=tasks
KEYCLOAK_CLIENT_ID=api
KEYCLOAK_CLIENT_SECRET=... # copy it from keycloak client credentials tab
```

For testing purpose, you can add an user to using keycloak and then try the `/login` endpoint for authentication to see whether it works fine or not.

It will start the server at port `8086`, and then you can perform any of the following requests:

- `POST /login`: Login with Keycloak user credentials
- `GET /tasks`: Get list of tasks
- `GET /tasks/:id`: Get a task with ID `id`
- `POST /tasks`: Create a new task
- `PATCH /tasks/:id`: Edit a task with ID `id` by providing `title`, `description`, `is_completed`
- `DELETE /tasks/:id`: Delete a task with ID `id`
- `GET /search/tasks?q=query`: Search tasks with title `query`

Here is an OpenAPI documentation of this API: [Swagger API Doc](https://app.swaggerhub.com/apis-docs/ARUPJANA7365_1/tasks-api/1.0.0)
