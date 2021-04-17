package postgres

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/novalwardhana/user-management-sso/global/constant"
	log "github.com/sirupsen/logrus"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func DBMasterRead() *gorm.DB {
	db := CreateConnection(os.Getenv(constant.ENVDBMaster))
	return db
}

func DBMasterWrite() *gorm.DB {
	db := CreateConnection(os.Getenv(constant.ENVDBMaster))
	return db
}

func CreateConnection(uri string) *gorm.DB {
	db, err := gorm.Open(os.Getenv(constant.ENVDBDriver), uri)

	if err != nil {
		log.Error(fmt.Sprintf("Cannot connect to db: %s", uri))
		i := 0
		for {
			i++
			reconnect, err := gorm.Open(os.Getenv(constant.ENVDBDriver), uri)
			if err != nil {
				log.Error(fmt.Sprintf("Cannot connect to db: %s", uri))
				time.Sleep(3 * time.Second)
				continue
			}
			log.Info(fmt.Sprintf("Successfully connected to db: %s", uri))
			return reconnect
		}
	}

	log.Info(fmt.Sprintf("Successfully connected to db: %s", uri))
	maxOpen, err := strconv.Atoi(constant.ENVMaxOpenConnection)
	if err != nil {
		maxOpen = 10
	}
	db.DB().SetMaxOpenConns(maxOpen)

	maxIdle, err := strconv.Atoi(constant.ENVMaxIdleConnection)
	if err != nil {
		maxIdle = 5
	}
	db.DB().SetMaxIdleConns(maxIdle)

	dbLog, err := strconv.ParseBool(os.Getenv(constant.ENVDBLog))
	if err != nil {
		dbLog = false
	}
	db.LogMode(dbLog)

	return db
}

func CloseConnection(db *gorm.DB) {
	if db != nil {
		db.Close()
		db = nil
	}
}
