package database

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"skframe/pkg/config"
	"skframe/pkg/console"
)

// DB 对象
var DB *gorm.DB
var SQLDB *sql.DB //设置数据库相关参数的谁会用到

func Connect(dbConfig gorm.Dialector, _logger gormLogger.Interface) {
	var err error
	DB, err = gorm.Open(dbConfig, &gorm.Config{
		Logger:         _logger,
		NamingStrategy: schema.NamingStrategy{SingularTable: true}, //禁用自动加s
	})
	if err != nil {
		console.Warning(err.Error())
	}
	// 获取底层的 sqlDB
	SQLDB, err = DB.DB()
	if err != nil {
		console.Warning(err.Error())
	}
}

func MasterOrSlave(master bool) (db *gorm.DB) { //使用主还是从
	switch master {
	case true:
		db = DB
	case false:
		db = DB
	}
	return
}

func CurrentDatabase() string { //当前数据库
	return DB.Migrator().CurrentDatabase()
}

func DeleteAllTables() error {
	var err error
	switch config.Get("database.connection") {
	case "mysql":
		err = deleteMySQLTables()
	case "sqlite":
		err = deleteAllSqliteTables()
	default:
		panic(errors.New("database connection not supported"))
	}

	return err
}
func deleteAllSqliteTables() error {
	tables := []string{}
	err := DB.Select(&tables, "SELECT name FROM sqlite_master WHERE type='table'").Error
	if err != nil {
		return err
	}
	// 删除所有表
	for _, table := range tables {
		err := DB.Migrator().DropTable(table)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteMySQLTables() error {
	dbname := CurrentDatabase()
	tables := []string{}
	// 读取所有数据表
	err := DB.Table("information_schema.tables").
		Where("table_schema = ?", dbname).
		Pluck("table_name", &tables).
		Error
	if err != nil {
		return err
	}
	// 暂时关闭外键检测
	DB.Exec("SET foreign_key_checks = 0;")
	// 删除所有表
	for _, table := range tables {
		err := DB.Migrator().DropTable(table)
		if err != nil {
			return err
		}
	}
	// 开启 MySQL 外键检测
	DB.Exec("SET foreign_key_checks = 1;")
	return nil
}
