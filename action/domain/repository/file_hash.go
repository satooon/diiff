package repository

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/satooon/diiff/action/domain/entity"
)

// FileHash repository
type FileHash struct {
	DB *gorm.DB
}

// DeleteByFilePath is record delete
func (f *FileHash) DeleteByFilePath(filePath ...string) error {
	return f.DB.Delete(entity.FileHash{}, fmt.Sprintf("file_path in ('%s')", strings.Join(filePath, "','"))).Error
}

// Create is new record create
func (f *FileHash) Create(filePath string, hash string) error {
	return f.DB.Create(&entity.FileHash{FilePath: filePath, Hash: hash}).Error
}

// Save is record update
func (f *FileHash) Save(filePath string, hash string) error {
	return f.DB.Save(&entity.FileHash{FilePath: filePath, Hash: hash}).Error
}
