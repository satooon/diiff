package entity

import "time"

// FileHash entity
type FileHash struct {
	FilePath  string `gorm:"primary_key"`
	Hash      string `gorm:"unique_index; not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
