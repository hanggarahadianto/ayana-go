package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteProject(c *gin.Context) {
	projectID := c.Param("id")
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var project models.Project
	if err := tx.Where("id = ?", projectID).First(&project).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Project not found", "error": err.Error()})
		return
	}

	// Delete Goods related to CashFlow
	if err := tx.Where("cash_flow_id IN (SELECT id FROM cash_flows WHERE project_id = ?)", projectID).Delete(&models.Goods{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to delete related Goods", "error": err.Error()})
		return
	}

	// Delete CashFlow records
	if err := tx.Where("project_id = ?", projectID).Delete(&models.CashFlow{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to delete related CashFlow records", "error": err.Error()})
		return
	}

	// Delete Workers and Materials related to WeeklyProgress
	if err := tx.Where("weekly_progress_id_worker IN (SELECT id FROM weekly_progresses WHERE project_id = ?)", projectID).Delete(&models.Worker{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to delete related Workers", "error": err.Error()})
		return
	}

	if err := tx.Where("weekly_progress_id_material IN (SELECT id FROM weekly_progresses WHERE project_id = ?)", projectID).Delete(&models.Material{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to delete related Materials", "error": err.Error()})
		return
	}

	// Delete WeeklyProgress records
	if err := tx.Where("project_id = ?", projectID).Delete(&models.WeeklyProgress{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to delete related WeeklyProgress records", "error": err.Error()})
		return
	}

	// Delete the Project itself
	if err := tx.Delete(&project).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to delete project", "error": err.Error()})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Project deleted successfully"})
}
