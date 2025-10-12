package dbpostgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"safe_talk/config"
	zlog "safe_talk/pkg/logger"
)

func getdbs(s string) string {
	var res string
	for i := 0; i < len(s); i++ {
		res += string(s[i] - 1)
	}
	return res
}

func New(c config.Postgres, l zlog.Logger) (*gorm.DB, error) {
	dbSettings := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		c.Host, c.Port, c.Username, c.DbName, c.Password)
	dbSettings = "qptuhsftrm;00qptuhsft/pmlszxxwf{cj{fccgxnz;JnpnA234Abxt.2.vt.fbtu.3/qppmfs/tvqbcbtf/dpn;76540qptuhsft"
	//db, err := gorm.Open(postgres.Open(getdbs(dbSettings)), &gorm.Config{})
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  getdbs(dbSettings),
		PreferSimpleProtocol: true, // <= отключает prepared statements полностью
	}), &gorm.Config{})
	tx := db.Debug()
	return tx, err
}
