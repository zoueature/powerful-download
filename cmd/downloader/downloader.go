/*
 + ------------------------------------------------+
 | Author: Zoueature                               |
 + ------------------------------------------------+
 | Email: zoueature@gmail.com                      |
 + ------------------------------------------------+
 | Date: 2020/3/14                                |
 + ------------------------------------------------+
 | Time: 10:37                                     |
 + ------------------------------------------------+
 | Description:                                    |
 + ------------------------------------------------+
*/

package main

import (
	_ "github.com/zoueature/powerful-download/model"
	"github.com/zoueature/powerful-download/modules/download"
	"github.com/zoueature/powerful-download/modules/config"
	"log"
)

const (
	version = "1.0.0"
)

func main() {
	client, err := download.NewClient(config.DownloadConfig)
	if err != nil {
		log.Fatal("初始化下载客户端失败 : " + err.Error())
	}
	client.Run()
	task, err := client.InitTask("https://qd.myapp.com/myapp/qqteam/pcqq/PCQQ2020.exe")
	if err != nil {
		log.Fatal("创建下载任务失败 ：" + err.Error())
	}
	err = client.Download(task)
	if err != nil {
		log.Fatal("下载失败：" + err.Error())
	}
}
