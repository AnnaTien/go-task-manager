# Go Task Manager

This is a simple task management application built with Go, using PostgreSQL as the database. It provides a REST API to manage tasks and a command-line interface (CLI) to interact with the API.

---

## System Requirements

Make sure you have the following tools installed on your machine:

* **Go** (version 1.23 or higher)
* **Docker**
* **Docker Compose**
* **cURL** (or a similar tool for API testing)

---

## Project Configuration

To run the project, you need to create a `.env` file to store sensitive environment variables.

1. Create a new file named **`.env`** at the project's root directory.
2. Add the following content to the file, replacing `your_user` and `your_password` with your credentials:

    ```text
    DB_HOST=db
    DB_USER=your_user
    DB_PASSWORD=your_password
    DB_NAME=task_manager_db
    DB_PORT=5432
    ```

---

## Getting Started

The project can be run in two ways: using Docker Compose or by running the CLI directly on your machine.

### Running with Docker Compose (Recommended)

This method will start both the Go application and the PostgreSQL database in separate containers.

1. Run the following command in the project's root directory to build and start the services:

    ```bash
    docker compose up --build
    ```

2. The API application will be running on port `8080`.

### Running the CLI Directly

You can build and run the Go application directly while still using the database managed by Docker.

1. Ensure the database container is running:

    ```bash
    docker compose up -d db
    ```

2. Build the CLI executable:

    ```bash
    go build -o go-task-manager ./cmd/
    ```

3. Run the API server:

    ```bash
    ./go-task-manager api
    ```

4. To connect successfully, ensure your `config/config.yaml` file uses `localhost` as the database host.

---

## Using the API

With the API server running, you can use the following `cURL` commands to interact with it.

* **Add a new task**

    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"name":"study unit test"}' http://localhost:8080/tasks
    ```

* **Get all tasks**

    ```bash
    curl http://localhost:8080/tasks
    ```

* **Search for a task**

    ```bash
    curl http://localhost:8080/tasks/search?q=study
    ```

* **Update a task**

    ```bash
    curl -X PUT -H "Content-Type: application/json" -d '{"name":"studyed unit test", "completed":true}' http://localhost:8080/tasks/1
    ```

* **Delete a task**

    ```bash
    curl -X DELETE http://localhost:8080/tasks/1
    ```

---

## Using the CLI

You can run the following CLI commands to interact with the API.

* **Add a new task**

    ```bash
    ./go-task-manager add "study Docker"
    ```

* **List all tasks**

    ```bash
    ./go-task-manager list
    ```

* **Delete a task**

    ```bash
    ./go-task-manager delete 1
    ```
