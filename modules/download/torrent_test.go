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
)

func TestBDecode(t *testing.T)  {
	f, _ := os.Open("./test.torrent")
	content, _ := ioutil.ReadAll(f)
	result, info, _ := bDecode(content)
	hashClient := sha1.New()
	hashClient.Write(info)
	infoHash := hashClient.Sum(nil)
	torrentInfo := result.(map[string]interface{})
	if torrentInfo == nil {
		log.Fatal("error")
	}
	announceList := torrentInfo["announce-list"].([]interface{})
	urlList1 := announceList[0].([]interface{})
	urlStr := urlList1[0].(string)
	base, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal("url error")
	}
	peeID := "sadaldajfqpowiedpqsa"
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
	responseInfo, _, _ := bDecode(responseStr)
	fmt.Println( responseInfo)

}


//http://t.t789.me:2710/announce?peer_id=jiujhuijkloiuytrewqa&info_hash=0a7dabeb053390089b3342087a4ede8a36d8716d&port=6881&left=0&downloaded=0&uploaded=0&compact=1


//http://t.t789.vip:2710/announce?info_hash=4ED89E23C6E13C171D2ADEBC5B5F7F6F5646E696&peer_id=0987ujikjyuhgsft11qw&port=8080&uploaded=0&downloaded=0&left=1263502296