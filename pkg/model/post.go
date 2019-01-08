package model

import "time"

// Post 文章
// status: 0,1 草稿，发布
type Post struct {
	ID              int64     `xorm:"pk autoincr INT(11)" json:"id"`
	UserId          int64     `xorm:"index INT(11)" json:"user_id"`
	Status          int       `xorm:"TINYINT(11)" json:"status"`
	Title           string    `xorm:"VARCHAR(255)" json:"title"`
	Path            string    `xorm:"unique VARCHAR(255)" json:"path"`
	MarkdownContent string    `xorm:"LONGTEXT" json:"markdown_content,omitempty"`
	CreateTime      time.Time `xorm:"index DATETIME" json:"create_time"`
	UpdateTime      time.Time `xorm:"DATETIME" json:"update_time"`
}
