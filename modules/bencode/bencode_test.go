/*+-----------------------------+
 *| Author: Zoueature           |
 *+-----------------------------+
 *| Email: zoueature@gmail.com  |
 *+-----------------------------+
 */
package bencode

import (
	"io/ioutil"
	"log"
	"testing"
	"fmt"
)

type BtInfo struct {
	Announce     string   `bencode:"announce"`
	AnnounceList [][]string `bencode:"announce-list"`
	Info         Info   `bencode:"info"`
}

type Info struct {
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
}

func TestBDecode(t *testing.T) {
	btContent, err := ioutil.ReadFile("test.torrent")
	if err != nil {
		log.Fatal(err.Error())
		t.Error(err)
	}
	container := new(BtInfo)
	err = BDecode(btContent, container)
	if err != nil {
		log.Fatal(err.Error())
		t.Error(err)
	}
	fmt.Println(container)
}
