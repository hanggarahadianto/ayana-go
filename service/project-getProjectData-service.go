package service

import (
	"ayana/db"
	"ayana/models"
	"ayana/utils/helper"
)

type ProjectFilterParams struct {
	CompanyID  string
	Pagination helper.Pagination
	DateFilter helper.DateFilter
	Search     string
}

func GetProjects(params ProjectFilterParams) ([]models.Project, int64, int64, error) {
	var (
		projects      []models.Project
		totalProject  int64
		filteredTotal int64
	)

	// Hitung semua project tanpa filter
	if err := db.DB.Model(&models.Project{}).Count(&totalProject).Error; err != nil {
		return nil, 0, 0, err
	}

	// Base query dengan filter company_id
	query := db.DB.Model(&models.Project{}).Where("company_id = ?", params.CompanyID)

	// Filter pencarian nama
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where(`
			project_name ILIKE ? OR 
			project_leader ILIKE ? OR 
			investor ILIKE ?`,
			searchPattern, searchPattern, searchPattern)
	}
	// Filter berdasarkan rentang tanggal ProjectStart
	if params.DateFilter.StartDate != nil && params.DateFilter.EndDate != nil {
		query = query.Where("project_start BETWEEN ? AND ?", params.DateFilter.StartDate, params.DateFilter.EndDate)
	} else if params.DateFilter.StartDate != nil {
		query = query.Where("project_start >= ?", params.DateFilter.StartDate)
	} else if params.DateFilter.EndDate != nil {
		query = query.Where("project_start <= ?", params.DateFilter.EndDate)
	}

	// Hitung total setelah filter (tanpa pagination)
	if err := query.Count(&filteredTotal).Error; err != nil {
		return nil, 0, 0, err
	}

	// Ambil data dengan pagination
	if err := query.
		Order("created_at desc, updated_at desc").
		Limit(params.Pagination.Limit).
		Offset((params.Pagination.Page - 1) * params.Pagination.Limit).
		Find(&projects).Error; err != nil {
		return nil, 0, 0, err
	}

	return projects, totalProject, filteredTotal, nil
}
