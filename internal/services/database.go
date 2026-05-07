package services

import (
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"sync"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type DatabaseService struct {
	DB          *gorm.DB
	logger      *slog.Logger
	dbPath      string
	startupOnce sync.Once
	startupErr  error
}

func NewDatabaseService(databasePath string, dbName string) *DatabaseService {
	p := filepath.Join(databasePath, dbName)
	return &DatabaseService{
		dbPath: p,
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

func (db *DatabaseService) WithLogger(logger *slog.Logger) *DatabaseService {
	db.logger = logger
	return db
}

func (db *DatabaseService) Startup() error {
	db.startupOnce.Do(func() {
		db.logger.Debug("Starting database service", "dbPath", db.dbPath)

		ndb, err := gorm.Open(sqlite.Open(db.dbPath), &gorm.Config{})
		if err != nil {
			db.startupErr = fmt.Errorf("open sqlite at %s: %w", db.dbPath, err)
			db.logger.Error("Failed to start database service", "error", db.startupErr.Error())
			return
		}

		db.DB = ndb
		db.logger.Info("Database service started successfully", "dbPath", db.dbPath)
	})
	return db.startupErr
}

func (db *DatabaseService) Close() error {
	if db.DB == nil {
		db.logger.Debug("Database connection is not initialized, skipping close")
		return nil
	}
	db.logger.Debug("Closing database connection", "dbPath", db.dbPath)

	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Close()
	if err != nil {
		db.logger.Error("Failed to close database connection", "error", err.Error())
	} else {
		db.logger.Info("Database connection closed successfully", "dbPath", db.dbPath)
	}

	return err
}
