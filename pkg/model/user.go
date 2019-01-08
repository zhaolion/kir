package model

import "time"

// User model.
type User struct {
	ID            int64     `xorm:"pk autoincr INT(11)" json:"id"`
	Name          string    `xorm:"index VARCHAR(191)" json:"name"`
	Nick          string    `xorm:"VARCHAR(191)" json:"nick"`
	Password      string    `xorm:"VARCHAR(191)" json:"password"`
	Salt          string    `xorm:"VARCHAR(191)" json:"salt"`
	AvatarURL     string    `xorm:"VARCHAR(191)" json:"avatar_url"`
	Email         string    `xorm:"index VARCHAR(191)" json:"email"`
	EmailVerified bool      `xorm:"VARCHAR(191)" json:"email_verified"`
	IsAdmin       bool      `xorm:"VARCHAR(191)" json:"is_admin"`
	CreateTime    time.Time `xorm:"index DATETIME" json:"create_time"`
	UpdateTime    time.Time `xorm:"DATETIME" json:"update_time"`
}
