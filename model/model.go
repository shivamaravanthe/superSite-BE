package model

import "time"

type Users struct {
	ID        uint      `gorm:"primary_key;AUTO_INCREMENT;column:ID;type:uint;"`
	Email     string    `gorm:"unique;column:Email;type:varchar(128);"`
	Password  string    `gorm:"column:Password;type:varchar(128)"`
	CreatedAt time.Time `gorm:"column:CreatedAt;type:timestamp"`
	CreatedBy string    `gorm:"column:CreatedBy;type:varchar(128)"`
	UpdatedAt time.Time `gorm:"column:UpdatedAt;type:timestamp"`
	UpdatedBy string    `gorm:"column:UpdatedBy;type:varchar(128)"`
}

type PasswordStore struct {
	ID          uint      `gorm:"primary_key;AUTO_INCREMENT;column:ID;type:uint;"`
	Link        string    `gorm:"column:Link;type:blob"`
	UserName    string    `gorm:"column:UserName;type:blob"`
	Password    string    `gorm:"column:Password;type:blob"`
	Description string    `gorm:"column:Description;type:varchar(128)"`
	CreatedAt   time.Time `gorm:"column:CreatedAt;type:timestamp"`
	CreatedBy   string    `gorm:"column:CreatedBy;type:varchar(128)"`
	UpdatedAt   time.Time `gorm:"column:UpdatedAt;type:timestamp"`
	UpdatedBy   string    `gorm:"column:UpdatedBy;type:varchar(128)"`
}
