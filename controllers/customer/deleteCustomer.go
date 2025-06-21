package controllers

import (
	"ayana/db"
	"ayana/models"
	"ayana/service"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteCustomer(c *gin.Context) {
	id := c.Param("id")

	// Periksa apakah customer ada
	var customer models.Customer
	if err := db.DB.First(&customer, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	// Hapus dari database
	if err := db.DB.Delete(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
		return
	}

	// Hapus dari Typesense
	err := service.DeleteCustomerFromTypesense(context.Background(), id)
	if err != nil {
		// Tidak menghalangi response sukses, hanya log error
		log.Printf("Gagal menghapus dokumen dari Typesense untuk ID %s: %v", id, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}
