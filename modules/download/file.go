/*+-----------------------------+
 *| Author: Zoueature           |
 *+-----------------------------+
 *| Email: zoueature@gmail.com  |
 *+-----------------------------+
 */
package download

type downloadInfo struct {
	totalLength         int64  //文件总长度
	fileName            string //文件名
	IP                  string //服务器IP
	PORT                int    //服务器端口
	supportPartDownload bool   //是否支持分段下载
}
