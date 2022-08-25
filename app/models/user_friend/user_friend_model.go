//Package user_friend 模型
package user_friend

import (
	"gorm.io/gorm"
	"skframe/app/models"
	"skframe/pkg/database"
)

type UserFriend struct {
	models.BaseModel
	UserId    uint64 `json:"user_id"`
	FriendId  uint64 `json:"friend_id"`
	AliasName string `json:"alias_name"`
	models.CommonTimestampsField
}

func (userFriend *UserFriend) Create(tx *gorm.DB) error {
	if tx != nil {
		return tx.Create(userFriend).Error
	} else {
		return database.DB.Create(userFriend).Error
	}
}

func (userFriend *UserFriend) CreateBatch(batch []UserFriend, tx *gorm.DB) error {
	if tx != nil {
		return tx.Create(batch).Error
	} else {
		return database.DB.Create(batch).Error
	}
}

func (userFriend *UserFriend) Del(opt models.SqlOpt, tx *gorm.DB) error {
	if tx == nil {
		tx = database.DB
	}
	for key, val := range opt.Where {
		tx = tx.Where(key, val)
	}
	return tx.Delete(UserFriend{}).Error
}

func (userFriend *UserFriend) Update(data map[string]interface{}, where map[string]interface{}, tx *gorm.DB) error {
	if tx == nil {
		tx = database.DB
	}
	for key, val := range where {
		tx = tx.Where(key, val)
	}
	return tx.Model(UserFriend{}).Updates(data).Error
}
