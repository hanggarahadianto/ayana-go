package routes

import (
	transactionCategory "ayana/controllers/transactionCategory"
	middlewares "ayana/middlewares/auth"

	"github.com/gin-gonic/gin"
)

func SetupTransactionCategoryRouter(r *gin.Engine) {
	transactionController := r.Group("/transaction-category", middlewares.AuthMiddleware())
	{
		transactionController.POST("/post", transactionCategory.CreateTransactionCategory)
		transactionController.GET("/get", transactionCategory.GetTransactionCategory)
		transactionController.PUT("/edit/:id", transactionCategory.UpdateTransactionCategory)
		transactionController.DELETE("/delete/:id", transactionCategory.DeleteTransactionCategory)

	}
}
