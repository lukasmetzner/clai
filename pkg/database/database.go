package database

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lukasmetzner/clai/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitDB() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_DATABASE")
	dbSslMode := os.Getenv("DB_SSLMODE")

	if dbSslMode == "" {
		dbSslMode = "disable"
	}

	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSslMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get database connection:", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	automigrate := os.Getenv("DB_AUTOMIGRATE")

	if automigrate == "" {
		automigrate = "0"
	}

	if automigrate == "1" {
		if err := db.AutoMigrate(&models.Job{}); err != nil {
			log.Fatal("failed to migrate database:", err)
		}

		if err := db.AutoMigrate(&models.JobOutput{}); err != nil {
			log.Fatal("failed to migrate database:", err)
		}
	}

	DB = db
}

func AppendJobOutput(jobID uint, stdout *bytes.Buffer, stderr *bytes.Buffer) {
	var job models.Job

	// Select from jobs where id = jobID
	if err := DB.First(&job, jobID).Error; err != nil {
		log.Printf("Error fetching job %d: %s", jobID, err)
	}

	jobOutput := models.JobOutput{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if err := DB.Create(&jobOutput).Error; err != nil {
		log.Printf("Error creating job output: %s", err)
	}

	job.JobOutputID = &jobOutput.ID
	if err := DB.Save(&job).Error; err != nil {
		log.Printf("Error saving job object: %s", err)
	}
}
