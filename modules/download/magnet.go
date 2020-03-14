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

func newMagnet() downloader {
	return &MagnetDownloader{}
}

func (mc *MagnetDownloader) download(task *downTask) error {
	return nil
}

func (mc *MagnetDownloader) cancel(task *downTask) error {
	return nil
}

func (mc *MagnetDownloader) parseURLInfo(url string) (*downloadInfo, error) {
	return nil, nil
}
