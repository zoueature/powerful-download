/*
 + ------------------------------------------------+
 | Author: Zoueature                               |
 + ------------------------------------------------+
 | Email: zoueature@gmail.com                      |
 + ------------------------------------------------+
 | Date: 2020/3/14                                |
 + ------------------------------------------------+
 | Time: 13:41                                     |
 + ------------------------------------------------+
 | Description:                                    |
 + ------------------------------------------------+
*/

package download

import (
	"io/ioutil"
	"os"
	"strconv"
)

type torrentData struct {
	dataType byte
	data     interface{}
	child    *torrentData
}

const (
	bDecodeInt    = 'i'
	bDecodeList   = 'l'
	bDecodeHash   = 'd'
	bDecodeString = 's'
	bDecodeEnd    = 'e'
)

var test = "d13:announce-listll30:http://t.t789.co:2710/announceel30:http://t.t789.me:2710/announceel31:http://t.t789.vip:2710/announceee4:infod6:lengthi1263502296e4:name71:[电影天堂www.dytt89.com]局内人2：头号通缉BD中英双字.mp412:piece lengthi524288e6:pieces48200:"

type TorrentDownloader struct {
	Path string
	IP   string
	Port int
}

func isTorrent(url string) bool {
	return false
}

func newTorrent() downloader {
	return &TorrentDownloader{}
}

func (tc *TorrentDownloader) download(task *downTask) error {
	return nil
}

func (tc *TorrentDownloader) cancel(task *downTask) error {
	return nil
}

func (tc *TorrentDownloader) parseURLInfo(url string) (*downloadInfo, error) {
	fileReader, err := os.Open(url)
	if err != nil {
		return nil, err
	}
	torrentContent, err := ioutil.ReadAll(fileReader)
	if err != nil {
		return nil, err
	}
	//torrentContent := []byte(url)
	//bdecode
	var typeStack []byte
	var matchContainer [][]byte
	var strMatcher []byte  //存储匹配的字符串长度
	var numMatcher []byte  //存储匹配的数值
	var startNumMatch bool //标识是否开启数值匹配
	info := new(torrentData)
	for index, b := range torrentContent {
		var nowType byte
		if len(typeStack) > 0 {
			nowType = typeStack[len(typeStack)-1]
		}
		if b == bDecodeHash || b == bDecodeList {
			typeStack = append(typeStack, b)
		} else if b == bDecodeInt {
			startNumMatch = true
		} else if b >= '0' && b <= '9' {
			if startNumMatch {
				numMatcher = append(numMatcher, b)
			} else {
				strMatcher = append(strMatcher, b)
			}
		} else if b == ':' {
			//字符串长度值匹配结束
			strLenStr := string(strMatcher)
			strLen, err := strconv.Atoi(strLenStr)
			if err != nil {
				return nil, err
			}
			str := torrentContent[index+1:index+strLen]
			matchContainer = append(matchContainer, str)
			strMatcher = []byte{}
		} else if b == bDecodeEnd {
			if startNumMatch {
				//数值匹配
				matchContainer = append(matchContainer, numMatcher)
				startNumMatch = false
				numMatcher = []byte{}
				continue
			}
		}
	}
	return nil, nil
}
