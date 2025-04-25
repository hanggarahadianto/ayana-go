package helper

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Pagination struct untuk memuat informasi halaman dan batas
type Pagination struct {
	Page   int
	Limit  int
	Offset int
}

// GetPagination dari query string di gin.Context
func GetPagination(c *gin.Context) Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	return Pagination{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

// ValidatePagination untuk validasi agar tidak ada limit atau page negatif
func ValidatePagination(p Pagination, c *gin.Context) bool {
	if p.Page < 1 || p.Limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return false
	}
	return true
}
