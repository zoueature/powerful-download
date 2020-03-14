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
	return nil, nil
}
