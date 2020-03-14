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
	"flag"
	"github.com/zoueature/powerful-download/modules/config"
	"github.com/zoueature/powerful-download/modules/download"
	"log"
)

const (
	version = "1.0.0"
)

func main() {
	var url string
	flag.StringVar(&url, "--url", "", "下载URL")
	downloader, err := download.NewParser(url, config.DownloadConfig)
	if err != nil {
		return
	}
	err = downloader.Parse()
	if err != nil {
		return
	}
	err = downloader.Download()
	if err != nil {
		log.Println(err.Error())
	}
}
