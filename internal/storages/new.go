package storages

import (
	"fmt"

	"github.com/Arup3201/gotask/internal/config"
	"github.com/Arup3201/gotask/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(config *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		config.Database.Host, config.Database.Port, config.Database.User,
		config.Database.Password, config.Database.Name, config.Database.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.User{}, &models.Task{})

	sql, _ := db.DB()
	sql.SetMaxOpenConns(config.Database.MaxConns)
	sql.SetMaxIdleConns(config.Database.MaxIdle)

	return db, nil
}
