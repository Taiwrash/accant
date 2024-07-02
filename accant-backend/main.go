package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("expenses.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Expense{})

	r := gin.Default()

	r.POST("/api/expenses", addExpense)
	r.GET("/api/expenses", getExpenses)
	r.GET("/api/sync-request-network", syncRequestNetworkMain)

	r.Run(":8080")
}

type Expense struct {
	gorm.Model
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

func addExpense(c *gin.Context) {
	var expense Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	db.Create(&expense)
	c.JSON(200, expense)
}

func getExpenses(c *gin.Context) {
	var expenses []Expense
	db.Find(&expenses)
	c.JSON(200, expenses)
}

func syncRequestNetworkMain(c *gin.Context) {
	syncRequestNetwork(c)
	c.JSON(200, gin.H{"message": "Sync completed"})
}
