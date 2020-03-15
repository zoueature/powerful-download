/*
 + ------------------------------------------------+
 | Author: Zoueature                               |
 + ------------------------------------------------+
 | Email: zoueature@gmail.com                      |
 + ------------------------------------------------+
 | Date: 2020/3/14                                |
 + ------------------------------------------------+
 | Time: 14:22                                     |
 + ------------------------------------------------+
 | Description:                                    |
 + ------------------------------------------------+
*/

package config

import (
	parenv "github.com/zoueature/par_env"
	"github.com/zoueature/powerful-download/modules/download"
	"log"
	"strconv"
)

const configFilePath = "./config/download.config"

var DownloadConfig download.DownConfig

func init() {
	parenv.EnvInit(configFilePath)
	threadNumStr := parenv.Get("downloadThreadNumber", "0")
	threadNum, err := strconv.Atoi(threadNumStr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	downLoadPath := parenv.Get("downloadPath", "")
	maxSpeedStr := parenv.Get("maxSpeed", "0")
	maxSpeed, err := strconv.Atoi(maxSpeedStr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	timeoutStr := parenv.Get("connectTimeOut", "0")
	connectTimeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	DownloadConfig = download.NewConfig(threadNum, downLoadPath, maxSpeed, connectTimeout)
}
