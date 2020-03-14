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

func (mc *MagnetDownloader) Parse() error {
	return nil
}

func (mc *MagnetDownloader) Download() error {
	return nil
}
