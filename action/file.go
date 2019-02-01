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
	Map() map[string]*fileInfo
}

type filePath struct {
	fileMap map[string]*fileInfo
}

type fileInfo struct {
	os.FileInfo
	path string
}

// NewFilePath return *file
func NewFilePath() FilePath {
	return &filePath{
		fileMap: map[string]*fileInfo{},
	}
}

// Scan is directory scan
func (f *filePath) Scan(path string) error {
	return filepath.Walk(path, f.walk)
}

func (f *filePath) walk(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return err
	}
	fi := &fileInfo{info, path}
	f.fileMap[fi.Path()] = fi
	return err
}

// Map file map
func (f *filePath) Map() map[string]*fileInfo {
	return f.fileMap
}

// Hash return string
func (fi *fileInfo) Hash() string {
	c := sha256.Sum256([]byte(fmt.Sprintf("%v%v%v", fi.Name(), fi.Size(), fi.ModTime())))
	return hex.EncodeToString(c[:])
}

func (fi *fileInfo) Path() string {
	return fi.path
}
