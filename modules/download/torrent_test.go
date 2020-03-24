/*+-----------------------------+
 *| Author: Zoueature           |
 *+-----------------------------+
 *| Email: zoueature@gmail.com  |
 *+-----------------------------+
 */
package download

import (
	"testing"
	"os"
	"io/ioutil"
	"crypto/sha1"
	"log"
	"strconv"
	"net/url"
	"net/http"
	"fmt"
	"net"
	"time"
)

func TestBDecode(t *testing.T)  {
	f, _ := os.Open("./test.torrent")
	content, _ := ioutil.ReadAll(f)
	result, info, err := bDecode(content)
	if err != nil {
		log.Fatal(err.Error())
	}
	hashClient := sha1.New()
	hashClient.Write(info)
	infoHash := hashClient.Sum(nil)
	torrentInfo := result.(map[string]interface{})
	if torrentInfo == nil {
		log.Fatal("error")
	}
	announceList := torrentInfo["announce-list"].([]interface{})
	for k := 0; k < 3; k ++ {
		urlList1 := announceList[k].([]interface{})
		urlStr := urlList1[0].(string)
		base, err := url.Parse(urlStr)
		if err != nil {
			log.Fatal("url error")
		}
		peeID := "xunlei-2715772xs1232"
		params := url.Values{
			"info_hash":  []string{string(infoHash[:])},
			"peer_id":    []string{peeID},
			"port":       []string{strconv.Itoa(int(6881))},
			"uploaded":   []string{"0"},
			"downloaded": []string{"0"},
			"compact":    []string{"1"},
			"left":       []string{strconv.Itoa(1263502296)},
		}
		base.RawQuery = params.Encode()
		requestUrl := base.String()
		httpClient := http.Client{}
		request, _ := http.NewRequest(http.MethodGet, requestUrl, nil)
		response , err := httpClient.Do(request)
		responseStr, _ := ioutil.ReadAll(response.Body)
		responseInfo, _, err := bDecode(responseStr)
		if err != nil {
			fmt.Println(err.Error())
		}
		re := responseInfo.(map[string]interface{})
		peersStr, ok  := re["peers"].(string)
		var peers []peer
		if ok {
			fmt.Println([]byte(peersStr))
			peers, err = peerDecode([]byte(peersStr))
			if err != nil {
				log.Fatal(err.Error())
			}
		}

		for _, p := range peers {
			host := net.JoinHostPort(p.ip.String(), strconv.Itoa(p.port))
			fmt.Println(host)
			conn, err := net.DialTimeout("tcp", host, 1 * time.Second)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			//hand shake
			protocol := "BitTorrent protocol"
			handShake := make([]byte, 1, len(protocol) + 49)
			handShake[0] = byte(len(protocol))
			handShake = append(handShake, []byte(protocol)...)
			handShake = append(handShake, make([]byte, 8)...)
			handShake = append(handShake, infoHash...)
			handShake = append(handShake, peeID...)
			fmt.Println(string(handShake))
			size, err := conn.Write(handShake)
			if err != nil {
				log.Println(err.Error())
			}
			fmt.Println("发送" + strconv.Itoa(size) + " 字节")
			buf := make([]byte, 100)
			conn.Read(buf)
			if err != nil {
				log.Println(err.Error())
			}
			fmt.Println(buf)
		}
	}
}


//http://t.t789.me:2710/announce?peer_id=jiujhuijkloiuytrewqa&info_hash=0a7dabeb053390089b3342087a4ede8a36d8716d&port=6881&left=0&downloaded=0&uploaded=0&compact=1


//http://t.t789.vip:2710/announce?info_hash=4ED89E23C6E13C171D2ADEBC5B5F7F6F5646E696&peer_id=0987ujikjyuhgsft11qw&port=8080&uploaded=0&downloaded=0&left=1263502296


//182.47.84.189