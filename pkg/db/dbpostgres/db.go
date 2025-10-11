package dbpostgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"safe_talk/config"
	zlog "safe_talk/pkg/logger"
)

func New(c config.Postgres, l zlog.Logger) (*gorm.DB, error) {
	dbSettings := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		c.Host, c.Port, c.Username, c.DbName, c.Password)
	dbSettings = "postgresql://postgres.olkrywwvezbizebbfwmy:Imom@123@aws-1-us-east-2.pooler.supabase.com:6543/postgres"
	db, err := gorm.Open(postgres.Open(dbSettings), &gorm.Config{})
	tx := db.Debug()
	tx.Callback()
	return tx, err
}
