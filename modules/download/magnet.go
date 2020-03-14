/*
 + ------------------------------------------------+
 | Author: Zoueature                               |
 + ------------------------------------------------+
 | Email: zoueature@gmail.com                      |
 + ------------------------------------------------+
 | Date: 2020/3/14                                |
 + ------------------------------------------------+
 | Time: 13:34                                     |
 + ------------------------------------------------+
 | Description:                                    |
 + ------------------------------------------------+
*/

package download

type MagnetDownloader struct {
	URL  string
	IP   string
	PORT int
}

func isMagnet(url string) bool {
	return false
}

func newMagnet()  {

}

func (mc *MagnetDownloader) Download(task downTask) error {
	return nil
}
