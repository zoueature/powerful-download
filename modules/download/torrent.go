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
	"crypto/sha1"
	"errors"
	"fmt"
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
	bDecodeInt      = 'i'
	bDecodeList     = 'l'
	bDecodeHash     = 'd'
	bDecodeString   = 's'
	bDecodeEnd      = 'e'
	bDecodeFormated = 'f'
)

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
	torrentInfo, infoContent, err := bDecode(torrentContent)
	if err != nil {
		return nil, err
	}
	_ = torrentInfo.(map[string]interface{})
	info := &downloadInfo{
		//todo format torrent info to download info
	}
	h := sha1.New()
	size, err := h.Write(infoContent)
	if err != nil {
		return nil, err
	}
	if size != len(infoContent){
		return nil, errors.New("check info hash error, write size not match")
	}
	infoHash := fmt.Sprintf("%x", h.Sum(nil))
	info.infoHash = infoHash
	return info, nil
}

func bDecode(torrentContent []byte) (interface{}, []byte, error) {
	//bdecode
	var typeStack []byte
	var matchContainer []interface{}
	var strMatcher []byte  //存储匹配的字符串长度
	var numMatcher []byte  //存储匹配的数值
	var startNumMatch bool //标识是否开启数值匹配
	var firstType byte
	var matchInfoLevel, infoStartIndex, infoEndIndex int
	for i := 0; i < len(torrentContent); i++ {
		b := torrentContent[i]
		if b == bDecodeHash || b == bDecodeList {
			if i == 0 {
				firstType = b
			}
			typeStack = append(typeStack, b)
			if len(matchContainer) > 0 && matchInfoLevel == 0 {
				v, ok := matchContainer[len(matchContainer)-1].(string)
				if ok && v == "info"{
					matchInfoLevel ++
					infoStartIndex = i
				}
			} else if  matchInfoLevel > 0 {
				matchInfoLevel ++
			}
		} else if b == bDecodeInt {
			startNumMatch = true
			if  matchInfoLevel > 0 {
				matchInfoLevel ++
			}
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
				return nil, nil, err
			}
			str := torrentContent[i+1:i+1+strLen]
			i += strLen
			matchContainer = append(matchContainer, string(str))
			strMatcher = append(strMatcher[0:0])
			typeStack = append(typeStack, bDecodeString)
		} else if b == bDecodeEnd {
			if matchInfoLevel == 1 {
				infoEndIndex = i + 1
				matchInfoLevel --
			} else if matchInfoLevel > 0 {
				matchInfoLevel --
			}
			if startNumMatch {
				//数值匹配
				matchContainer = append(matchContainer, string(numMatcher))
				startNumMatch = false
				numMatcher = append(numMatcher[0:0])
				typeStack = append(typeStack, bDecodeInt)
				continue
			}
			tmp := make([]interface{}, 0)
			var nowType byte
			typeLen := len(typeStack)
			var j int
			for j = 0; j < typeLen; j ++ {
				nowType = typeStack[len(typeStack)-j-1]
				if nowType == bDecodeFormated || nowType == bDecodeInt || nowType == bDecodeString {
					tmp = append(tmp, matchContainer[len(matchContainer)-j-1])
				} else {
					break
				}
			}
			if len(tmp) == 0 {
				return nil, nil, errors.New("format data error")
			}

			matchContainer = append(matchContainer[:len(matchContainer)-j])
			typeStack = append(typeStack[:len(typeStack)-j-1])
			var data interface{}
			if nowType == bDecodeList {
				data = tmp
			} else if nowType == bDecodeHash {
				l := len(tmp)
				if l % 2 != 0 {
					return nil, nil, errors.New("format map error, item num error ")
				}
				m := make(map[string]interface{})
				var key string
				for k := l; k > 0; k -- {
					index := k - 1
					if k % 2 == 0 {
						var ok bool
						key, ok = tmp[index].(string)
						if !ok {
							return nil, nil, errors.New("format map error, trans to key string error ")
						}
					} else {
						m[key] = tmp[index]
					}
				}
				data = m
			}
			matchContainer = append(matchContainer, data)
			typeStack = append(typeStack, bDecodeFormated)
		}
	}
	if len(matchContainer) < 0 {
		return nil, nil, errors.New("error,  bdecode empty")
	}
	info := torrentContent[infoStartIndex:infoEndIndex]
	if firstType == bDecodeHash {
		return matchContainer[0], info, nil
	}
	return matchContainer, info, nil
}

