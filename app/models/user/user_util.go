package user

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"skframe/app/models"
	"skframe/pkg/database"
)

func First(opt models.SqlOpt, tx *gorm.DB, share bool) (user User) {
	if tx == nil {
		tx = database.DB
	} else if share == true {
		tx = tx.Clauses(clause.Locking{Strength: "SHARE", Table: clause.Table{Name: clause.CurrentTable}}) //共享锁，可读不可以修改
	} else {
		tx = tx.Clauses(clause.Locking{Strength: "UPDATE"}) //排他锁,禁止其他读，当前会修改
	}
	for key, val := range opt.Where {
		tx = tx.Where(key, val)
	}
	if opt.Order != "" {
		tx = tx.Order(opt.Order)
	}
	tx.Select(opt.Field).First(&user)
	return
}

func Find(opt models.SqlOpt, tx *gorm.DB, share bool) (user []User) {
	if tx == nil {
		tx = database.DB
	} else if share == true {
		tx = tx.Clauses(clause.Locking{Strength: "SHARE", Table: clause.Table{Name: clause.CurrentTable}}) //共享锁，可读不可以修改
	} else {
		tx = tx.Clauses(clause.Locking{Strength: "UPDATE"}) //排他锁,禁止其他读，当前会修改
	}
	for key, val := range opt.Where {
		tx = tx.Where(key, val)
	}
	if opt.Order != "" {
		tx = tx.Order(opt.Order)
	}
	if opt.Limit > 0 {
		tx = tx.Limit(opt.Limit)
	}
	if opt.Offset > 0 {
		tx = tx.Limit(opt.Offset)
	}
	tx.Select(opt.Field).Find(&user)
	return
}

func IsExist(where map[string]interface{}) bool {
	tx := database.DB
	for key, val := range where {
		tx = tx.Where(key, val)
	}
	var count int64
	tx.Model(User{}).Count(&count)
	return count > 0
}

func Count(opt models.SqlOpt, tx *gorm.DB, share bool) (count int64) {
	if tx == nil {
		tx = database.DB
	} else if share == true {
		tx = tx.Clauses(clause.Locking{Strength: "SHARE", Table: clause.Table{Name: clause.CurrentTable}}) //共享锁，可读不可以修改
	} else {
		tx = tx.Clauses(clause.Locking{Strength: "UPDATE"}) //排他锁,禁止其他读，当前会修改
	}
	for key, val := range opt.Where {
		tx = tx.Where(key, val)
	}
	tx.Model(User{}).Count(&count)
	return
}

func Paginate(opt models.SqlOpt, pages, size int) (count int64, list []User) {
	dbTx := database.DB.Model(User{}).Select(opt.Field)
	for key, val := range opt.Where {
		dbTx = dbTx.Where(key, val)
	}
	dbTx.Count(&count)
	dbTx.Offset(pages*size - pages).Limit(size).Find(&list)
	return
}
