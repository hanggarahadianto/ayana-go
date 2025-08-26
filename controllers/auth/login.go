// controllers/login.go
package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	utilsAuth "ayana/utils/auth"
	utilsEnv "ayana/utils/env"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type CompanyUserLite struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CompanyResp struct {
	ID          uuid.UUID         `json:"id"`
	Title       string            `json:"title"`
	CompanyCode string            `json:"company_code"`
	Color       string            `json:"color"`
	HasCustomer bool              `json:"has_customer"`
	HasProject  bool              `json:"has_project"`
	HasProduct  bool              `json:"has_product"`
	IsRetail    bool              `json:"is_retail"`
	Users       []CompanyUserLite `json:"users"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type AuthUserResp struct {
	ID        uuid.UUID     `json:"id"`
	Username  string        `json:"username"`
	Password  string        `json:"password"`
	Role      string        `json:"role"`
	Companies []CompanyResp `json:"companies"` // slice, bukan pointer
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type AuthDataResp struct {
	User  AuthUserResp `json:"user"`
	Token string       `json:"token"`
}

type AuthEnvelope struct {
	Status  bool         `json:"status"`
	Message string       `json:"message"`
	Data    AuthDataResp `json:"data"`
}

type LoginData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Login(c *gin.Context) {
	var loginData LoginData
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	var user models.User
	if err := db.DB.First(&user, "username = ?", loginData.Username).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Username not found",
		})
		return
	}

	if err := utilsAuth.VerifiedPassword(user.Password, loginData.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "Wrong password",
		})
		return
	}

	// --- Ambil companies sesuai role ---
	var companies []models.Company
	if user.Role == "superadmin" {
		// Semua company + preload relasi untuk dapatkan users (via user_companies)
		if err := db.DB.
			Preload("UserCompanies").
			Preload("UserCompanies.User").
			Find(&companies).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to fetch companies",
			})
			return
		}
	} else {
		// Company milik user ini
		if err := db.DB.
			Model(&models.Company{}).
			Joins("JOIN user_companies uc ON uc.company_id = companies.id").
			Where("uc.user_id = ?", user.ID).
			Preload("UserCompanies").
			Preload("UserCompanies.User").
			Find(&companies).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to fetch user companies",
			})
			return
		}
	}

	// --- Map companies -> []CompanyResp (users: ICompanyUser[]) ---
	companiesResp := make([]CompanyResp, 0) // 👈 penting: jangan biarin nil
	for _, comp := range companies {
		usersLite := make([]CompanyUserLite, 0)
		for _, uc := range comp.UserCompanies {
			usersLite = append(usersLite, CompanyUserLite{
				ID:        uc.User.ID,
				Username:  uc.User.Username,
				Role:      uc.User.Role,
				CreatedAt: uc.User.CreatedAt,
				UpdatedAt: uc.User.UpdatedAt,
			})
		}

		companiesResp = append(companiesResp, CompanyResp{
			ID:          comp.ID,
			Title:       comp.Title,
			CompanyCode: comp.CompanyCode,
			Color:       comp.Color,
			HasCustomer: comp.HasCustomer,
			HasProject:  comp.HasProject,
			HasProduct:  comp.HasProduct,
			IsRetail:    comp.IsRetail,
			Users:       usersLite,
			CreatedAt:   comp.CreatedAt,
			UpdatedAt:   comp.UpdatedAt,
		})
	}

	// --- JWT ---
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	config, _ := utilsEnv.LoadConfig(".")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to generate token",
		})
		return
	}

	// Cookie HttpOnly
	c.SetCookie("token", tokenString, 3600*24, "/", "", false, true)

	// --- Bentuk response sesuai IAuthResponse ---
	resp := AuthEnvelope{
		Status:  true,
		Message: "Login success",
		Data: AuthDataResp{
			Token: tokenString,
			User: AuthUserResp{
				ID:        user.ID,
				Username:  user.Username,
				Password:  user.Password, // hashed
				Role:      user.Role,
				Companies: companiesResp, // 👈 sudah pasti [] bukan null
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			},
		},
	}

	c.JSON(http.StatusOK, resp)
}
