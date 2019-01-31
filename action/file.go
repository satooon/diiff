package action

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

// FilePath interface
type FilePath interface {
	Scan(path string) error
}

type filePath struct {
	files []*fileInfo
}

type fileInfo struct {
	os.FileInfo
}

// NewFilePath return *file
func NewFilePath() FilePath {
	return &filePath{
		files: []*fileInfo{},
	}
}

// Scan is directory scan
func (f *filePath) Scan(path string) error {
	return filepath.Walk(path, f.walk)
}

func (f *filePath) walk(path string, info os.FileInfo, err error) error {
	f.files = append(f.files, &fileInfo{info})
	return err
}

// Hash return string
func (fi *fileInfo) Hash() string {
	c := sha256.Sum256([]byte(fmt.Sprintf("%v%v%v", fi.Name(), fi.Size(), fi.ModTime())))
	return hex.EncodeToString(c[:])
}
