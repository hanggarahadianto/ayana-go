package service

import (
	"ayana/db"
	"ayana/models"
	"ayana/utils/helper"
)

type ProjectFilterParams struct {
	CompanyID  string
	Pagination helper.Pagination
	Search     string // Jika ingin menambahkan pencarian di masa depan
}

func GetProjects(params ProjectFilterParams) ([]models.Project, int64, int64, error) {
	var (
		projects      []models.Project
		totalProject  int64
		filteredTotal int64
	)

	// Total seluruh project tanpa filter
	if err := db.DB.Model(&models.Project{}).Count(&totalProject).Error; err != nil {
		return nil, 0, 0, err
	}

	// Query dengan filter company_id dan search
	query := db.DB.Model(&models.Project{}).Where("company_id = ?", params.CompanyID)

	if params.Search != "" {
		query = query.Where("name ILIKE ?", "%"+params.Search+"%")
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
