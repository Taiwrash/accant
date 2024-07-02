package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	requestNetworkAPIURL = "https://api.request.network/requests" // Replace with actual API endpoint
	requestNetworkAPIKey = "Key"                                  // Replace with your actual API key
)

type RequestNetworkTransaction struct {
	ID          string  `json:"id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

func syncRequestNetwork(c *gin.Context) {
	// Make a request to the Request Network API
	client := &http.Client{}
	req, err := http.NewRequest("GET", requestNetworkAPIURL, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Add("Authorization", "Bearer "+requestNetworkAPIKey)

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch data from Request Network"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read response body"})
		return
	}

	var transactions []RequestNetworkTransaction
	err = json.Unmarshal(body, &transactions)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to parse response"})
		return
	}

	// Process and store transactions
	for _, transaction := range transactions {
		expense := Expense{
			Amount:      transaction.Amount,
			Description: transaction.Description,
			Date:        transaction.Date,
		}

		result := db.Create(&expense)
		if result.Error != nil {
			fmt.Printf("Failed to store transaction %s: %v\n", transaction.ID, result.Error)
		}
	}

	c.JSON(200, gin.H{"message": "Sync completed", "transactions_processed": len(transactions)})
}
