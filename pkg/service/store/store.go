package store

import (
	"github.com/go-xorm/xorm"
)

type DB struct {
	engine *xorm.Engine
}
