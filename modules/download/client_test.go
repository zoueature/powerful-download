/*+-----------------------------+
 *| Author: Zoueature           |
 *+-----------------------------+
 *| Email: zoueature@gmail.com  |
 *+-----------------------------+
 */
package download

import (
	"testing"
	"github.com/zoueature/powerful-download/modules/msg"
	"time"
	"fmt"
)

func TestNewClient(t *testing.T) {
	testHttp := "https://qd.myapp.com/myapp/qqteam/pcqq/PCQQ2020.exe"
	config := NewConfig(0, "", 0, 0)
	client, err := NewClient(config)
	if err != nil {
		t.Error("创建客户端失败: " + err.Error())
	}
	task, err := client.InitTask(testHttp)
	client.Download(task)
}

func TestClient_Download(t *testing.T) {
	config := NewConfig(0, "", 0, 0)
	client, err := NewClient(config)
	if err != nil {
		t.Error("创建客户端失败: " + err.Error())
	}
	client.Run()
	time.Sleep(1 * time.Second)
	task := &msg.OperatorMsg{}
	task.MsgType = msg.OperateMsgNewDownload
	task.MsgContent = "https://qd.myapp.com/myapp/qqteam/pcqq/PCQQ2020.exe"
	fmt.Println("put task")
	client.operatorChan <- task
	time.Sleep(1000 * time.Second)
}