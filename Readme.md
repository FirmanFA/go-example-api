# Project Name

A sample backend RESTful API project written in GoLang using Gin framework and MySQL.

## Overview

This project provides APIs for managing projects and their budgets. It includes endpoints for retrieving projects, creating new projects, updating existing projects, and deleting projects.

## Setup

1. Install GoLang on your system if not already installed.
2. Clone this repository to your local machine.
3. Install required dependencies using `go mod tidy`.
4. Set up MySQL and create a database named `recordings`.
5. Set environment variables `DBUSER` and `DBPASS` with your MySQL username and password.
6. Run the project using `go run main.go`.

## API Endpoints

### Get All Projects

- **URL:** `/projects`
- **Method:** `GET`
- **Description:** Retrieves all projects along with their budget details.
- **Response:** Array of project objects.

### Get Project by ID

- **URL:** `/projects/:id`
- **Method:** `GET`
- **Description:** Retrieves a specific project by its ID.
- **Response:** Project object.

### Create New Project

- **URL:** `/projects`
- **Method:** `POST`
- **Description:** Creates a new project with the provided details.
- **Request Body:** JSON object with project details (title, leader, budget).
- **Response:** Created project object.

### Update Project

- **URL:** `/projects/:id`
- **Method:** `PUT`
- **Description:** Updates an existing project by its ID.
- **Request Body:** JSON object with updated project details (title, leader, budget).
- **Response:** Updated project object.

### Delete Project

- **URL:** `/projects/:id`
- **Method:** `DELETE`
- **Description:** Deletes a project along with its associated budget by its ID.
- **Response:** Success message.

## Dependencies

- `github.com/gin-gonic/gin`: Web framework for building APIs in GoLang.
- `github.com/go-sql-driver/mysql`: MySQL driver for GoLang.

## Database Schema

- **project:**

  - id (PK)
  - title
  - leader

- **project_budget:**
  - budget_value
  - down_payment
  - deadline
  - project_id (FK referencing project.id)

## Author

[Your Name](https://github.com/yourusername)

## License

This project is licensed under the [MIT License](LICENSE).
