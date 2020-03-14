/*
 + ------------------------------------------------+
 | Author: Zoueature                               |
 + ------------------------------------------------+
 | Email: zoueature@gmail.com                      |
 + ------------------------------------------------+
 | Date: 2020/3/14                                |
 + ------------------------------------------------+
 | Time: 13:33                                     |
 + ------------------------------------------------+
 | Description:                                    |
 + ------------------------------------------------+
*/

package download

import (
	"log"
	"regexp"
)

type HttpDownloader struct {
}

func isHttp(url string) bool {
	//过滤完前后空格的url
	reg, err := regexp.Compile(`^http(s)?://[\s\S]]+`)
	if err != nil {
		log.Println("download url error, URL:" + url + ", error: " + err.Error())
		return false
	}
	match := reg.Match([]byte(url))
	return match
}

func newHTTP() Downloader {
	return &HttpDownloader{}
}

func (hc *HttpDownloader) Download(task downTask) error {

}

func (hc *HttpDownloader) Cancel(task downTask) error {

}