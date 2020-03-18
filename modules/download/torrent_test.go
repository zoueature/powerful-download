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
	"fmt"
)

func TestBDecode(t *testing.T)  {
	f, _ := os.Open("./test.torrent")
	content, _ := ioutil.ReadAll(f)
	info, _ := bDecode(content)

	tInfo := info.(map[string]interface{})
	torrentInfoInterface := tInfo["info"]
	torrentInfo := torrentInfoInterface.(map[string]interface{})
	piecesInterface := torrentInfo["pieces"]
	piece := piecesInterface.(string)
	pieces := []byte(piece)
	result := make([]string, 0)
	start, end := 0, 20
	l := len(pieces)
	for {
		if start >= l {
			break
		}
		tmp := pieces[start:end]
		result = append(result, fmt.Sprintf("%x", tmp))
		start = end
		end += 20
	}
	fmt.Println(result)

}


//http://t.t789.me:2710/announce?peer_id=jiujhuijkloiuytrewqa&info_hash=0a7dabeb053390089b3342087a4ede8a36d8716d&port=6881&left=0&downloaded=0&uploaded=0&compact=1