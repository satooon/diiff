package action

import (
	"github.com/jinzhu/gorm"
)

// DB interface
type DB interface {
	Open(dialect string, args ...interface{}) error
	DB() *gorm.DB
}

type db struct {
	_db *gorm.DB
}

// NewDB return DB
func NewDB() DB {
	return &db{}
}

func (d *db) Open(dialect string, args ...interface{}) error {
	_db, err := gorm.Open(dialect, args...)
	if err != nil {
		return err
	}
	d._db = _db
	return nil
}

func (d *db) DB() *gorm.DB {
	return d._db
}
