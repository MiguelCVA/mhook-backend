package migrations

import (
	"github.com/MiguelCVA/mhook-backend/internal/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	db.AutoMigrate(&models.User{}, &models.Project{}, &models.User{}, &models.Session{})
}
