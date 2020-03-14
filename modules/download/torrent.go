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

type TorrentDownloader struct {
	Path string
	IP   string
	Port int
}

func isTorrent(url string) bool {
	return false
}

func (tc *TorrentDownloader) Parse() error {
	return nil
}

func (tc *TorrentDownloader) Download(task downTask) error {
	return nil
}
