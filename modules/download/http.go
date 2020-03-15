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
	"net/http"
	"regexp"
	"strings"
)

type HttpDownloader struct {
	client http.Client
}

func isHttp(url string) bool {
	//过滤完前后空格的url
	reg, err := regexp.Compile(`^http(s)?://[\s\S]+`)
	if err != nil {
		log.Println("download url error, URL:" + url + ", error: " + err.Error())
		return false
	}
	match := reg.Match([]byte(url))
	return match
}

func newHTTP() downloader {
	return &HttpDownloader{
		client: http.Client{},
	}
}

func (hc *HttpDownloader) cancel(task *downTask) error {
	return nil
}

func (hc *HttpDownloader) parseURLInfo(url string) (*downloadInfo, error) {
	request, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := hc.client.Do(request)
	if err != nil {
		return nil, err
	}
	allPath := strings.Split(url, `/`)
	fileName := allPath[len(allPath)-1]
	info := &downloadInfo{
		fileName: fileName,
	}
	contentLength := response.ContentLength
	info.totalLength = contentLength
	acceptRange := response.Header.Get("Accept-Ranges")
	response.Body.Close()
	if acceptRange == "bytes" {
		info.supportPartDownload = true
	}
	return info, nil
}

func (hc *HttpDownloader) download(task *downTask) error {
	task.status = taskStatusDownloading
	task.listenSignal()
	for i := 0; i < task.threadNum; i ++ {
		go func() {
			content, err := http.Get(task.URL)
			if err != nil {

			}
		}()
	}
	return nil
}
