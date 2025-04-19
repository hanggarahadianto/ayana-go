package routes

import (
	transactionCategory "ayana/controllers/transactionCategory"

	"github.com/gin-gonic/gin"
)

func SetupTransactionCategoryRouter(r *gin.Engine) {
	transactionController := r.Group("/transaction-category")
	{
		transactionController.POST("/post", transactionCategory.CreateTransactionCategory)
		transactionController.GET("/get", transactionCategory.GetTransactionCategory)

	}
}
