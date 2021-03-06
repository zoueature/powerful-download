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
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
)

type torrentData struct {
	mainAnnounce string
	announces    []string
	infoHash     []byte
	fileName     string
	length       int64
	pieces       [][20]byte
}

type peer struct {
	ip   net.IP
	port int
}

const (
	//bdecode
	bDecodeInt      = 'i'
	bDecodeList     = 'l'
	bDecodeHash     = 'd'
	bDecodeString   = 's'
	bDecodeEnd      = 'e'
	bDecodeFormated = 'f'

	//peer
	peerPerLength = 6 // 4ip+2port
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
	td := torrentData{}
	ti := torrentInfo.(map[string]interface{})
	shaClient := sha1.New()
	shaClient.Write(infoContent)
	td.infoHash = shaClient.Sum(nil)
	mainAnnounce, ok := ti["announce"]
	if ok {
		td.mainAnnounce = mainAnnounce.(string)
	}
	announces, ok := ti["announce-list"]
	if ok {
		list, ok := announces.([]interface{})
		if ok {
			announceList := make([]string, len(list))
			for _, v := range list {
				vl, ok := v.([]interface{})
				if ok {
					if vl[0] != nil {
						announce, ok := vl[0].(string)
						if ok {
							announceList = append(announceList, announce)
						}
					}
					continue
				}
				vs, ok := v.(string)
				if ok && vs != "" {
					announceList = append(announceList, vs)
				}
			}
			td.announces = announceList
		}
	}
	nameI, ok := ti["name"]
	if ok {
		name, ok := nameI.(string)
		if ok {
			td.fileName = name
		}
	}
	lengthI, ok := ti["length"]
	if ok {
		lengthStr, ok := lengthI.(string)
		if ok {
			length, err := strconv.ParseInt(lengthStr, 10, 64)
			if err != nil {
				return nil, err
			}
			td.length = length
		}
	}
	info := &downloadInfo{
		//todo format torrent info to download info
	}
	h := sha1.New()
	size, err := h.Write(infoContent)
	if err != nil {
		return nil, err
	}
	if size != len(infoContent) {
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
				if ok && v == "info" {
					matchInfoLevel++
					infoStartIndex = i
				}
			} else if matchInfoLevel > 0 {
				matchInfoLevel++
			}
		} else if b == bDecodeInt {
			startNumMatch = true
			if matchInfoLevel > 0 {
				matchInfoLevel++
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
			str := torrentContent[i+1 : i+1+strLen]
			i += strLen
			matchContainer = append(matchContainer, string(str))
			strMatcher = append(strMatcher[0:0])
			typeStack = append(typeStack, bDecodeString)
		} else if b == bDecodeEnd {
			if matchInfoLevel == 1 {
				infoEndIndex = i + 1
				matchInfoLevel--
			} else if matchInfoLevel > 0 {
				matchInfoLevel--
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
			for j = 0; j < typeLen; j++ {
				nowType = typeStack[len(typeStack)-j-1]
				if nowType == bDecodeFormated || nowType == bDecodeInt || nowType == bDecodeString {
					tmp = append(tmp, matchContainer[len(matchContainer)-j-1])
				} else {
					break
				}
			}
			//if len(tmp) == 0 {
			//	return nil, nil, errors.New("format data error")
			//}

			matchContainer = append(matchContainer[:len(matchContainer)-j])
			typeStack = append(typeStack[:len(typeStack)-j-1])
			var data interface{}
			if nowType == bDecodeList {
				data = tmp
			} else if nowType == bDecodeHash {
				l := len(tmp)
				if l%2 != 0 {
					return nil, nil, errors.New("format map error, item num error ")
				}
				m := make(map[string]interface{})
				var key string
				for k := l; k > 0; k-- {
					index := k - 1
					if k%2 == 0 {
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

func peerDecode(peers []byte) ([]peer, error) {
	const peerLen = 6
	if len(peers)%peerLen != 0 {
		return nil, errors.New("error byte num")
	}
	ls := make([]peer, 0, len(peers)%peerLen)
	for i := 0; i < len(peers)/peerLen; i++ {
		startIndex := i * peerPerLength
		endIndex := (i + 1) * peerPerLength
		p := peers[startIndex:endIndex]
		tmp := peer{}
		ip := make([]byte, 4)
		//tmp.ip = net.IP(p[:4])
		copy(ip, p[:4])
		tmp.ip = ip
		tmp.port = int(binary.BigEndian.Uint16(p[4:]))
		ls = append(ls, tmp)
	}
	return ls, nil
}
