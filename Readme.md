# Go Backend Sample Project

A sample backend RESTful API project written in GoLang using Gin framework and MySQL.

## Before Running
Set environment variables `DBUSER` and `DBPASS` with MySQL username and password.

## Dependencies

- `github.com/gin-gonic/gin`: Web framework for building APIs in GoLang.
- `github.com/go-sql-driver/mysql`: MySQL driver for GoLang.

## API Endpoints

| Endpoint             | Method | Description                                         | Request Body                  | Response                 |
|----------------------|--------|-----------------------------------------------------|-------------------------------|--------------------------|
| `/projects`          | GET    | Retrieves all projects with budget details.         | N/A                           | Array of project objects |
| `/projects/:id`      | GET    | Retrieves a specific project by ID.                 | N/A                           | Project object           |
| `/projects`          | POST   | Creates a new project.                              | JSON (title, leader, budget) | Created project object   |
| `/projects/:id`      | PUT    | Updates an existing project by ID.                  | JSON (title, leader, budget) | Updated project object   |
| `/projects/:id`      | DELETE | Deletes a project by ID along with its budget.      | N/A                           | Success message          |
