/*+-----------------------------+
 *| Author: Zoueature           |
 *+-----------------------------+
 *| Email: zoueature@gmail.com  |
 *+-----------------------------+
 */
package download

import (
	"fmt"
	"testing"
)

func TestBDecode(t *testing.T)  {
	var test = "d13:announce-listll30:http://t.t789.co:2710/announceel30:http://t.t789.me:2710/announceel31:http://t.t789.vip:2710/announceee4:infod6:lengthi1263502296e4:name71:[电影天堂www.dytt89.com]局内人2：头号通缉BD中英双字.mp412:piece lengthi524288e6:pieces1:de"
	info, err := bDecode([]byte(test))
	fmt.Println(info, err)
}