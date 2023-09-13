package models

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type UserProfile struct {
	Id         int64     `json:"id" gorm:"Id"`
	NickName   string    `json:"nick_name" gorm:"nick_name"`
	Pwd        string    `json:"pwd" gorm:"pwd"`
	Email      string    `json:"email" gorm:"email"`
	CreateTime time.Time `json:"create_time" gorm:"create_time"`
	UpdateTime time.Time `json:"update_time" gorm:"update_time"`
	UserId     string    `json:"user_id" gorm:"user_id"`
	HeadIcon   int       `json:"head_icon" gorm:"head_icon"`
	Level      int       `json:"level" gorm:"level"`
}

func (u *UserProfile) TableName() string {
	return "user.user_profile"
}
func (u *UserProfile) GetByEmail(db *gorm.DB, email string) error {
	err := db.Table(u.TableName()).Where("email = ?", email).Find(u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return err
}
func (u *UserProfile) GetByUserIds(db *gorm.DB, userIds []string) (list []UserProfile, err error) {
	if len(userIds) == 0 {
		return
	}
	err = db.Table(u.TableName()).Where("user_id in (?)", userIds).Find(&list).Error
	return
}
func (u *UserProfile) Create(db *gorm.DB) error {
	return db.Table(u.TableName()).Create(u).Error
}
