package model

// Tag 标签
type Tag struct {
	ID    int    `xorm:"pk autoincr INT(11)" json:"id"`
	Name  string `xorm:"unique VARCHAR(64)" json:"name"`
	Intro string `xorm:"VARCHAR(64)" json:"intro"`
}
