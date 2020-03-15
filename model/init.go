/*+-----------------------------+
 *| Author: Zoueature           |
 *+-----------------------------+
 *| Email: zoueature@gmail.com  |
 *+-----------------------------+
 */
package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"os"
	"errors"
)

const (
	//db connect info
	dbType          = "sqlite3"
	dbDataDirectory = "./data"
	dbDataFile      = dbDataDirectory + "/download.db"

	//tables
	taskDbTable = "t_task"
)

type connection struct {
	db *gorm.DB
}

var conn *connection

func init() {
	err := tryToMakeDataFile()
	if err != nil {
		log.Fatal("init database error : " + err.Error())
	}
	db, err := gorm.Open(dbType, dbDataFile)
	if err != nil {
		log.Fatal("connect database error : " + err.Error())
	}
	conn = &connection{
		db: db,
	}
}

func GetModel() *connection {
	return conn
}

func tryToMakeDataFile() error {
	info, err := os.Stat(dbDataFile)
	if err != nil {
		if os.IsNotExist(err) {
			//download.db文件不存在
			info, err := os.Stat(dbDataDirectory)
			if err != nil {
				if os.IsNotExist(err) {
					err := os.Mkdir(dbDataDirectory, 0755)
					if err != nil {
						return err
					}
				}
			} else {
				if !info.IsDir() {
					return errors.New("data目录被文件占用， 请处理后重试")
				}
			}
			_, err = os.Create(dbDataFile)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if info.IsDir() {
			return errors.New("数据文件被文件夹占用， 请处理后重启")
		}
	}
	return nil
}
