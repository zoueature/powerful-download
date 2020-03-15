/*+-----------------------------+
 *| Author: Zoueature           |
 *+-----------------------------+
 *| Email: zoueature@gmail.com  |
 *+-----------------------------+
 */
package msg

type OperatorMsg struct {
	Message
}

const (
	OperateMsgCancel      = 1  	//取消下载任务
	OperateMsgPause       = 2  	//暂停下载任务
	OperateMsgDelete      = 3	//删除下载任务
	OperateMsgContinue    = 4	//继续下载任务
	OperateMsgRedownload  = 5	//重新下载
	OperateMsgClientQuit  = 6	//客户端退出
	OperateMsgNewDownload = 7	//创建新的下载任务
)
