# Tasks API Documentation

## Endpoints

### GET /tasks

- **Description**: Retrieve a list of tasks.
- **Response**:
  - **200 OK**: Returns a list of tasks.
  - **500 Internal Server Error**: An error occurred while fetching tasks.

### POST /tasks

- **Description**: Create a new task.
- **Request Body**:
  ```json
  {
    "title": "Task Title",
    "description": "Task Description"
  }
  ```
- **Response**:
  - **201 Created**: Task created successfully.
  - **400 Bad Request**: Invalid request data.
  - **500 Internal Server Error**: An error occurred while creating the task.

### GET /tasks/{id}

- **Description**: Retrieve a specific task by ID.
- **Response**:
  - **200 OK**: Returns the task with the specified ID.
  - **404 Not Found**: Task not found.
  - **500 Internal Server Error**: An error occurred while fetching the task.

### PUT /tasks/{id}

- **Description**: Update a specific task by ID.
- **Request Body**:
  ```json
  {
    "title": "Updated Task Title",
    "description": "Updated Task Description"
  }
  ```
- **Response**:
  - **200 OK**: Task updated successfully.
  - **404 Not Found**: Task not found.
  - **500 Internal Server Error**: An error occurred while updating the task.

### DELETE /tasks/{id}

- **Description**: Delete a specific task by ID.
- **Response**:
  - **204 No Content**: Task deleted successfully.
  - **404 Not Found**: Task not found.
  - **500 Internal Server Error**: An error occurred while deleting the task.

## Environment Variables

- `PORT`: The port on which the API server will listen (default: :8888)
- `ENVIRONMENT`: The environment mode (default: production, development)
- `LOG_FORMAT`: The format of log messages (default: json, console)
- `LOG_LEVEL`: The log level (default: info, debug, warn, error)
- `LOG_TIME_FORMAT`: The time format for log messages (default: rfc3339, unix, etc.)
- `DB_HOST`: The hostname of the database server (default: db)
- `DB_PORT`: The port of the database server (default: 5432)
- `DB_USER`: The username for database authentication
- `DB_PASSWORD`: The password for database authentication
- `DB_NAME`: The name of the database
- `DB_SSL_MODE`: The SSL mode for database connections (default: disable)
- `CORS_ALLOWED_ORIGINS`: A comma-separated list of allowed origins for CORS (default: *)
- `CORS_ALLOWED_METHODS`: A comma-separated list of allowed HTTP methods for CORS (default: GET,POST,PUT,DELETE,OPTIONS)
- `CORS_ALLOWED_HEADERS`: A comma-separated list of allowed HTTP headers for CORS (default: Content-Type,Authorization)
- `CORS_EXPOSED_HEADERS`: A comma-separated list of exposed HTTP headers for CORS (default: Content-Type,Authorization)
- `CORS_ALLOW_CREDENTIALS`: Whether to allow credentials in CORS requests (default: true)
- `CORS_MAX_AGE`: The maximum age of a preflight request in seconds (default: 300)
