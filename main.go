package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type budgetModel struct {
	BudgetValue int64  `json:"budget_value"`
	DownPayment int64  `json:"down_payment"`
	Deadline    string `json:"deadline"`
}

type projectModel struct {
	ID     string      `json:"id"`
	Title  string      `json:"title"`
	Leader string      `json:"leader"`
	Budget budgetModel `json:"budget"`
}

var db *sql.DB

func main() {

	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	router := gin.Default()
	router.GET("/projects", getProjects)
	router.GET("/projects/:id", getProjectById)
	router.POST("/projects", postProjects)
	router.PUT("/project/:id", updateProject)
	router.DELETE("/project/:id", deleteProject)

	router.Run("localhost:8080")
}

func getProjects(c *gin.Context) {

	var projects []projectModel

	rows, err := db.Query("SELECT p.id, p.title, p.leader, pb.budget_value, pb.down_payment, pb.deadline FROM project p JOIN project_budget pb ON p.id = pb.project_id;")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var proj projectModel
		var projBudget budgetModel
		if err := rows.Scan(&proj.ID, &proj.Title, &proj.Leader, &projBudget.BudgetValue, &projBudget.DownPayment, &projBudget.Deadline); err != nil {
			return
		}
		proj.Budget = projBudget
		projects = append(projects, proj)
	}

	c.IndentedJSON(http.StatusOK, projects)
}

func getProjectById(c *gin.Context) {
	id := c.Param("id")

	var proj projectModel
	var projBudget budgetModel

	row := db.QueryRow("SELECT p.id, p.title, p.leader, pb.budget_value, pb.down_payment, pb.deadline FROM project p JOIN project_budget pb ON p.id = pb.project_id WHERE p.id = ?", id)
	if err := row.Scan(&proj.ID, &proj.Title, &proj.Leader, &projBudget.BudgetValue, &projBudget.DownPayment, &projBudget.Deadline); err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "project not found"})
			return
		}
		fmt.Println("Error scanning data:", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	proj.Budget = projBudget
	c.IndentedJSON(http.StatusOK, proj)
}

func postProjects(c *gin.Context) {
	var newProject projectModel

	//binding request to struct model
	if err := c.BindJSON(&newProject); err != nil {
		fmt.Println("Error binding json to struct:", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Error starting transaction:", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	// Insert query for the project table within the transaction
	projectQuery := "INSERT INTO project (title, leader) VALUES (?, ?)"
	projectResult, err := tx.Exec(projectQuery, newProject.Title, newProject.Leader)
	if err != nil {
		fmt.Println("Error inserting project to database:", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		// Rollback the transaction if there's an error
		tx.Rollback()
		return
	}

	projectID, err := projectResult.LastInsertId()

	if err != nil {
		fmt.Println("Error getting last inserted id:", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	// Insert query for the project_budget table within the transaction
	budgetQuery := "INSERT INTO project_budget (budget_value, down_payment, deadline, project_id) VALUES (?, ?, ?, ?)"
	_, err = tx.Exec(budgetQuery, newProject.Budget.BudgetValue, newProject.Budget.DownPayment, newProject.Budget.Deadline, projectID)
	if err != nil {
		fmt.Println("Error inserting into project_budget table:", err)

		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})

		tx.Rollback()
		return
	}

	// Commit the transaction if all insertions were successful
	err = tx.Commit()
	if err != nil {
		fmt.Println("Error committing transaction:", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	c.IndentedJSON(http.StatusCreated, newProject)

}

func updateProject(c *gin.Context) {

	id := c.Param("id")

	var newProject projectModel

	//binding request to struct model
	if err := c.BindJSON(&newProject); err != nil {
		fmt.Println("Error binding json to struct:", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Error starting transaction:", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	// Update query for the project table within the transaction
	projectQuery := "UPDATE project SET title = ?, leader = ? WHERE id = ?"
	updateProjectResult, err := tx.Exec(projectQuery, newProject.Title, newProject.Leader, id)
	if err != nil {
		fmt.Println("Error updating project table:", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		// Rollback the transaction if there's an error
		tx.Rollback()
		return
	}

	if rowAffected, _ := updateProjectResult.RowsAffected(); rowAffected == 0 {
		fmt.Println("No rows affected")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
		// Rollback the transaction if there's an error
		tx.Rollback()
		return
	}

	// Update query for the project_budget table within the transaction
	budgetQuery := "UPDATE project_budget SET budget_value = ?, down_payment = ?, deadline = ? WHERE project_id = ?"
	updateBudgetResult, err := tx.Exec(budgetQuery, newProject.Budget.BudgetValue, newProject.Budget.BudgetValue, newProject.Budget.Deadline, id)
	if err != nil {
		fmt.Println("Error updating project_budget table:", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		// Rollback the transaction if there's an error
		tx.Rollback()
		return
	}

	if rowAffected, _ := updateBudgetResult.RowsAffected(); rowAffected == 0 {
		fmt.Println("No rows affected")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
		// Rollback the transaction if there's an error
		tx.Rollback()
		return
	}

	// Commit the transaction if all updates were successful
	err = tx.Commit()
	if err != nil {
		fmt.Println("Error committing transaction:", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	c.IndentedJSON(http.StatusOK, newProject)

}

func deleteProject(c *gin.Context){

	id := c.Param("id")

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Error starting transaction:", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	// Delete query for the project_budget table within the transaction
	budgetQuery := "DELETE FROM project_budget WHERE project_id = ?"
	deleteBudgetResult, err := tx.Exec(budgetQuery, id)
	if err != nil {
		fmt.Println("Error deleting from project_budget table:", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		// Rollback the transaction if there's an error
		tx.Rollback()
		return
	}

	if rowAffected, _ := deleteBudgetResult.RowsAffected(); rowAffected == 0 {
		fmt.Println("No rows affected")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
		// Rollback the transaction if there's an error
		tx.Rollback()
		return
	}

	// Delete query for the project table within the transaction
	projectQuery := "DELETE FROM project WHERE id = ?"
	deleteProjectResult, err := tx.Exec(projectQuery, id)
	if err != nil {
		fmt.Println("Error deleting from project table:", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		// Rollback the transaction if there's an error
		tx.Rollback()
		return
	}

	if rowAffected, _ := deleteProjectResult.RowsAffected(); rowAffected == 0 {
		fmt.Println("No rows affected")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
		// Rollback the transaction if there's an error
		tx.Rollback()
		return
	}

	// Commit the transaction if all deletions were successful
	err = tx.Commit()
	if err != nil {
		fmt.Println("Error committing transaction:", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully deleted!"})

}
