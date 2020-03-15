/*+-----------------------------+
 *| Author: Zoueature           |
 *+-----------------------------+
 *| Email: zoueature@gmail.com  |
 *+-----------------------------+
 */
package msg

type Message struct {
	MsgType    int
	MsgContent interface{}
}

type TaskMsg struct {
	TaskID int
	Message
}

const (
	TaskMsgStop       = 1 //取消下载任务
	TaskMsgStart      = 2 //暂停下载任务
	TaskMsgRedownload = 3 //重新下载
	TaskMsgDelete     = 4 //删除任务
)
