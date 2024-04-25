package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

var projects = []projectModel{
	{
		ID:     "0",
		Title:  "Project RSUD",
		Leader: "Budi",
		Budget: budgetModel{
			BudgetValue: 19000,
			DownPayment: 900,
			Deadline:    "tomorrow"}},
	{ID: "1", Title: "Project Bupati", Leader: "Steven", Budget: budgetModel{BudgetValue: 19000, DownPayment: 900, Deadline: "tomorrow"}},
	{ID: "2", Title: "Project Dinas Kehutanan", Leader: "Gerard", Budget: budgetModel{BudgetValue: 19000, DownPayment: 900, Deadline: "tomorrowsss"}},
}

func main() {
	router := gin.Default()
	router.GET("/projects", getProjects)
	router.GET("/projects/:id", getProjectById)
	router.POST("/projects", postProjects)

	router.Run("localhost:8080")
}

func getProjects(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, projects)
}

func getProjectById(c *gin.Context) {
	id := c.Param("id")
	for _, p := range projects {
		if p.ID == id {
			c.IndentedJSON(http.StatusOK, p)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "project not found"})
}

func postProjects(c *gin.Context) {
	var newProject projectModel

	if err := c.BindJSON(&newProject); err != nil {
		return
	}

	projects = append(projects, newProject)
	c.IndentedJSON(http.StatusCreated, newProject)

}
