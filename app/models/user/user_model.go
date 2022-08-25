//Package user 模型
package user

import (
	"gorm.io/gorm"
	"skframe/app/models"
	"skframe/pkg/database"
)

type User struct {
	models.BaseModel
	Account  string `json:"account"`
	Password string `json:"password"`
	Name     string `json:"name"`
	models.CommonTimestampsField
}

func (user *User) Create(tx *gorm.DB) error {
	if tx != nil {
		return tx.Create(user).Error
	} else {
		return database.DB.Create(user).Error
	}
}

func (user *User) CreateBatch(batch []User, tx *gorm.DB) error {
	if tx != nil {
		return tx.Create(batch).Error
	} else {
		return database.DB.Create(batch).Error
	}
}

func (user *User) Del(opt models.SqlOpt, tx *gorm.DB) error {
	if tx == nil {
		tx = database.DB
	}
	for key, val := range opt.Where {
		tx = tx.Where(key, val)
	}
	return tx.Delete(User{}).Error
}

func (user *User) Update(data map[string]interface{}, where map[string]interface{}, tx *gorm.DB) error {
	if tx == nil {
		tx = database.DB
	}
	for key, val := range where {
		tx = tx.Where(key, val)
	}
	return tx.Model(User{}).Updates(data).Error
}
