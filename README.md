# Task API

[OpenAPI specification](https://app.swaggerhub.com/apis-docs/ARUPJANA7365_1/tasks-api/1.0.0)

## Set up

After cloning this repo, you can do the following to run the server.

```sh
go build cmd/api/main.go
./main.exe
```

It will start the server at port `8080`, and then you can perform any of the following requests:

- `GET /tasks`: Get list of tasks
- `GET /tasks/:id`: Get a task with ID `id`
- `POST /tasks`: Create a new task
- `PUT /tasks/:id`: Edit a task with ID `id` by providing `title` and `description`
- `PUT /tasks/:id/mark`: Mark a task as done/undone by providing `completed` as true/false
- `DELETE /tasks/:id`: Delete a task with ID `id`
- `GET /search/tasks?q=query`: Search tasks with title `query`