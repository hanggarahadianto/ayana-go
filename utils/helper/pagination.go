package helper

import (
	"net/http"
	"strconv"
	"time"

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

type DateFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
}

func GetDateFilter(c *gin.Context) (DateFilter, error) {
	layout := "2006-01-02" // format yyyy-mm-dd

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate *time.Time

	if startDateStr != "" {
		sd, err := time.Parse(layout, startDateStr)
		if err != nil {
			return DateFilter{}, err
		}
		startDate = &sd
	}

	if endDateStr != "" {
		ed, err := time.Parse(layout, endDateStr)
		if err != nil {
			return DateFilter{}, err
		}
		endDate = &ed
	}

	return DateFilter{
		StartDate: startDate,
		EndDate:   endDate,
	}, nil
}
