package db

import (
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

type McatDB struct {
	DB   *leveldb.DB
	Path string
}

var DefaultPath string = "./data/leveldb/"
var instance *McatDB

func init() {
	if _, err := os.Stat(DefaultPath); err != nil {
		os.Mkdir(DefaultPath, 0755)
	}
}
func NewDB(path string) (*McatDB, error) {
	if instance == nil {
		if _, err := os.Stat(path); err != nil {
			return nil, err
		}
		db, err := leveldb.OpenFile(path, nil)
		if err != nil {
			return nil, err
		}
		//defer db.Close()

		instance = new(McatDB)
		instance.DB = db
		instance.Path = path
	}
	return instance, nil
}

func (self *McatDB) Put(key, value []byte) error {
	return self.DB.Put(key, value, nil)
}

func (self *McatDB) Get(key []byte) ([]byte, error) {
	return self.DB.Get(key, nil)
}

func (self *McatDB) Delete(key []byte) error {
	return self.DB.Delete(key, nil)
}
