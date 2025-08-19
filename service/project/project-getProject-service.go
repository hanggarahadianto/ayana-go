package service

import (
	"ayana/db"
	lib "ayana/lib"
	"ayana/models"

	"gorm.io/gorm"
)

type ProjectFilterParams struct {
	CompanyID  string
	Pagination lib.Pagination
	DateFilter lib.DateFilter
	Search     string
}

func GetProjects(params ProjectFilterParams) ([]models.Project, int64, int64, error) {
	var (
		projects      []models.Project
		totalProject  int64
		filteredTotal int64
	)

	// Base query untuk filter company_id
	baseQuery := db.DB.Model(&models.Project{}).Where("company_id = ?", params.CompanyID)

	// Hitung total project (per company_id)
	if err := baseQuery.Count(&totalProject).Error; err != nil {
		return nil, 0, 0, err
	}

	// Copy query untuk filter tambahan
	query := baseQuery.Session(&gorm.Session{})

	// Filter pencarian nama
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where(`
			project_name ILIKE ? OR 
			project_leader ILIKE ? OR 
			investor ILIKE ?`,
			searchPattern, searchPattern, searchPattern)
	}

	orderBy := "project_start DESC"

	// Filter berdasarkan rentang tanggal ProjectStart
	if params.DateFilter.StartDate != nil && params.DateFilter.EndDate != nil {
		query = query.Where("project_start BETWEEN ? AND ?", params.DateFilter.StartDate, params.DateFilter.EndDate)
	} else if params.DateFilter.StartDate != nil {
		query = query.Where("project_start >= ?", params.DateFilter.StartDate)
	} else if params.DateFilter.EndDate != nil {
		query = query.Where("project_start <= ?", params.DateFilter.EndDate)
	}

	// Hitung total setelah filter
	if err := query.Count(&filteredTotal).Error; err != nil {
		return nil, 0, 0, err
	}

	// Ambil data dengan pagination
	if err := query.
		Order(orderBy).
		Limit(params.Pagination.Limit).
		Offset((params.Pagination.Page - 1) * params.Pagination.Limit).
		Find(&projects).Error; err != nil {
		return nil, 0, 0, err
	}

	return projects, totalProject, filteredTotal, nil
}
