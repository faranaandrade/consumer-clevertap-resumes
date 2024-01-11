package gorm

import (
	"fmt"
	"net/url"
	"time"

	"github.com/occmundial/go-common/logger"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	envProd = "prod"
)

type DBGorm struct {
	config *Setup
	log    logger.Logger
}

func NewDBGorm(cfg *Setup, log *logger.Log) *gorm.DB {
	d := DBGorm{
		config: cfg,
		log:    log,
	}
	dbGorm, err := d.Create()
	if err != nil {
		log.Fatal("gorm", "NewDBGorm", err)
	}
	return dbGorm
}

func (d *DBGorm) Create() (*gorm.DB, error) {
	connectionString := d.buildConnectionString()
	db, err := d.initializeDBSession(connectionString)
	if err != nil {
		return nil, err
	}
	d.setConnectionPool(db)

	return db, nil
}

func (d *DBGorm) buildConnectionString() string {
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d?ApplicationIntent=ReadOnly&database=%s&connection+timeout=%d&dial+timeout=%d",
		d.config.User, url.QueryEscape(d.config.Password), d.config.Server, d.config.Port, d.config.Database,
		d.config.ConnectTimeOut, d.config.ConnectTimeOut)
}

func (d *DBGorm) initializeDBSession(connectionString string) (*gorm.DB, error) {
	var logMode = gormlogger.Error
	if d.config.Environment == envProd {
		logMode = gormlogger.Silent
	}
	db, err := gorm.Open(sqlserver.Open(connectionString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			NoLowerCase: true, // skip the snake_casing of names
		},
		Logger: gormlogger.Default.LogMode(logMode)})

	if err != nil {
		d.log.Error("gorm", "initializeDBSession", err)
		return nil, err
	}
	return db, nil
}

func (d *DBGorm) setConnectionPool(db *gorm.DB) {
	sqlServer, err := db.DB()
	if err != nil {
		d.log.Error("gorm", "setConnectionPool", err)
	}
	sqlServer.SetMaxIdleConns(d.config.MaxIdleConnections)
	sqlServer.SetMaxOpenConns(d.config.MaxOpenConnections)
	sqlServer.SetConnMaxLifetime(time.Minute * time.Duration(d.config.ConnMaxLifetime))
}
