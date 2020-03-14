/*
 + ------------------------------------------------+
 | Author: Zoueature                               |
 + ------------------------------------------------+
 | Email: zoueature@gmail.com                      |
 + ------------------------------------------------+
 | Date: 2020/3/14                                |
 + ------------------------------------------------+
 | Time: 11:51                                     |
 + ------------------------------------------------+
 | Description:                                    |
 + ------------------------------------------------+
*/

package download

import (
	"errors"
	"strings"
)

type urlChecker func(string) bool
type downloadMaker func(string) Downloader

type Downloader interface {
	Parse() error
	Download() error
}

//下载配置
type DownConfig struct {
	ThreadNum    int    //线程数
	ToPath       string //下载路劲
	MaxSpeed     int    //最高下载速度
	Timeout      int    //超时时间
	BreakPointer int    //下载文件偏移量， 用于断点续传
}

const (
	httpClient    = "http"
	magnetClient  = "magnet"
	torrentClient = "torrent"
)

var (
	urlSupport = map[string]urlChecker{
		httpClient: isHttp,
	}
	downloaderSupport = map[string]downloadMaker{
		httpClient: newHTTP,
	}
)

func NewParser(url string) (Downloader, error) {
	url = strings.Trim(url, " ")
	var downloader Downloader
	var err error
	for key, checker := range urlSupport {
		if !checker(url) {
			continue
		}
		maker, ok := downloaderSupport[key]
		if !ok {
			err = errors.New("No support client ")
			break
		}
		downloader = maker(url)
	}
	return downloader, err
}
