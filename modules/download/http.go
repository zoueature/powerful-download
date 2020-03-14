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
	URL string
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

func newHTTP(url string) Downloader {
	return &HttpDownloader{
		URL: url,
	}
}

func (hc *HttpDownloader) Parse() error {
	return nil
}

func (hc *HttpDownloader) Download() error {
	return nil
}
