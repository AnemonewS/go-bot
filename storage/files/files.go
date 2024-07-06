package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"telegram-go/lib/e"
	"telegram-go/storage"
	"time"
)

type Storage struct {
	basePath string
}

const defaultPermission = 0774

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("Can't save page", err) }()

	filePath := filepath.Join(s.basePath, page.UserName) // For all OC, instead of path.Join(), where / != \ in windows;
	if err := os.MkdirAll(filePath, defaultPermission); err != nil {
		return err
	}
	fName, err := fileName(page)
	if err != nil {
		return err
	}
	filePath = filepath.Join(filePath, fName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	} // serialize page and write to file
	return nil
}

func (s Storage) ChoseRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("Can't pick random page", err) }()

	filePath := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(filePath)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}
	n := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(files))
	randFile := files[n]

	return s.decodePage(filepath.Join(filePath, randFile.Name()))
}

func (s Storage) Remove(page *storage.Page) error {
	fName, err := fileName(page)
	if err != nil {
		return e.WrapError("can't remove file", err)
	}
	path := filepath.Join(s.basePath, page.UserName, fName)
	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)
		return e.WrapError(msg, err)
	}
	return nil
}

func (s Storage) Exists(page *storage.Page) (bool, error) {
	fName, err := fileName(page)
	if err != nil {
		return false, e.WrapError("can't check file existence", err)
	}
	path := filepath.Join(s.basePath, page.UserName, fName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check file existence %s", path)
		return false, e.WrapError(msg, err)
	}
	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, e.WrapError("can't open file", err)
	}
	defer func() { _ = file.Close() }()
	var p storage.Page

	if err := gob.NewDecoder(file).Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func fileName(page *storage.Page) (string, error) {
	return page.Hash()
}
