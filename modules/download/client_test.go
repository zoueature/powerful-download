/*+-----------------------------+
 *| Author: Zoueature           |
 *+-----------------------------+
 *| Email: zoueature@gmail.com  |
 *+-----------------------------+
 */
package download

import (
	"testing"
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