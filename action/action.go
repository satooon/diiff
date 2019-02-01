package action

import (
	"fmt"
	"os"

	"github.com/satooon/diiff/action/domain/entity"
	"github.com/satooon/diiff/action/domain/repository"

	"github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3" // use sqlite
	"github.com/urfave/cli"
)

var (
	log = logrus.New()
)

// Action interface
type Action interface {
	Action(ctx *cli.Context) error
}

type action struct {
	debug bool
	db    DB
}

// NewAction return Action
func NewAction(debug bool) (Action, error) {
	db := NewDB()
	if err := db.Open("sqlite3", "diiff.db"); err != nil {
		return nil, err
	}
	db.DB().LogMode(debug)
	db.DB().AutoMigrate(&entity.FileHash{})

	a := &action{debug, db}
	log.Level = logrus.ErrorLevel
	if debug {
		log.Level = logrus.DebugLevel
	}
	log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
	log.SetReportCaller(debug)
	return a, nil
}

func (a *action) Action(ctx *cli.Context) error {
	log.Println("Args", ctx.Args())

	path, err := getPath(ctx.Args().Get(0))
	if err != nil {
		return err
	}
	log.Println("Path", path)

	fileNames := map[string]struct{}{}

	filePath := NewFilePath()
	if err = filePath.Scan(path); err != nil {
		return err
	}
	fileMap := filePath.Map()
	for k := range fileMap {
		fileNames[k] = struct{}{}
	}

	var fileHashes []*entity.FileHash
	if err = a.db.DB().Find(&fileHashes).Error; err != nil {
		return err
	}
	var fileHashMap = map[string]*entity.FileHash{}
	for _, v := range fileHashes {
		fileHashMap[v.FilePath] = v
		fileNames[v.FilePath] = struct{}{}
	}

	callFn, err := a.diff(fileNames, fileMap, fileHashMap)
	if err != nil {
		return err
	}
	a.db.DB().Begin()
	defer func() {
		if r := recover(); r != nil {
			a.db.DB().Rollback()
		}

	}()
	for _, v := range callFn {
		if err := v(); err != nil {
			a.db.DB().Rollback()
			return err
		}
	}
	a.db.DB().Commit()

	log.Println("Success")

	return nil
}

func (a *action) diff(fileNames map[string]struct{}, fileMap map[string]*fileInfo, fileHashMap map[string]*entity.FileHash) ([]func() error, error) {
	repo := repository.FileHash{DB: a.db.DB()}
	var callFn []func() error

	for k := range fileNames {
		log.Println("FileName", k)
		fi, ok := fileMap[k]
		if !ok {
			log.Println("DB Del", k)
			callFn = append(callFn, func(k string) func() error {
				return func() error {
					return repo.DeleteByFilePath(k)
				}
			}(k))
			continue
		}
		fh, ok := fileHashMap[k]
		if !ok {
			log.Println("DB Add", k)
			fmt.Printf("%s\n", fi.Path())
			callFn = append(callFn, func(f *fileInfo) func() error {
				return func() error {
					return repo.Create(f.Path(), f.Hash())
				}
			}(fi))
			continue
		}
		if fi.Hash() == fh.Hash {
			log.Println("Same", k)
			continue
		}
		log.Println("DB Upd", k)
		fmt.Printf("%s\n", fi.Path())
		callFn = append(callFn, func(fi *fileInfo) func() error {
			return func() error {
				return repo.Save(fi.Path(), fi.Hash())
			}
		}(fi))
	}

	return callFn, nil
}

func getPath(p string) (string, error) {
	if len(p) > 0 {
		return p, nil
	}
	return os.Getwd()
}
