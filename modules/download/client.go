/*
 + ------------------------------------------------+
 | Author: Zoueature                               |
 + ------------------------------------------------+
 | Email: zoueature@gmail.com                      |
 + ------------------------------------------------+
 | Date: 2020/3/14                                |
 + ------------------------------------------------+
 | Time: 11:51                                     |
 + ------------------------------------------------+
 | Description:                                    |
 + ------------------------------------------------+
*/

package download

import (
	"errors"
	"strings"
	"sync"
	"time"
)

type urlChecker func(string) bool
type downloadMaker func() Downloader

type Downloader interface {
	Download(task downTask) error
	Cancel(task downTask) error
}

type Client struct {
	DownConfig
	downloader    Downloader
	task          chan *downTask
	lock          sync.RWMutex
	runTaskNum    int
	communication map[int64]chan *msg
}

type msg struct {
	msgType int
}

const (
	httpClient    = "http"
	magnetClient  = "magnet"
	torrentClient = "torrent"
)

var (
	urlSupport = map[string]urlChecker{
		httpClient:    isHttp,
		magnetClient:  isMagnet,
		torrentClient: isTorrent,
	}
	downloaderSupport = map[string]downloadMaker{
		httpClient: newHTTP,
	}
	clients = map[string]*Client{}
)

//实例化一个下载客户端
func NewClient(url string, config DownConfig) (*Client, error) {
	url = strings.Trim(url, " ")
	var downloader Downloader
	var err error
	task := &downTask{
		taskID: time.Now().UnixNano(),
		URL:    url,
	}
	var clientType string
	for key, checker := range urlSupport {
		if !checker(url) {
			continue
		}
		client, ok := clients[key]
		if ok {
			client.PutTask(task)
			return client, nil
		}
		maker, ok := downloaderSupport[key]
		if !ok {
			err = errors.New("No support client ")
			break
		}
		downloader = maker()
		clientType = key
		break
	}
	if err != nil {
		return nil, err
	}
	//任务队列
	taskChan := make(chan *downTask, config.maxTask)
	//通信通道
	communication := make(map[int64]chan *msg, config.maxRunTask)
	communication[task.taskID] = make(chan *msg)
	client := &Client{
		downloader:    downloader,
		task:          taskChan,
		communication: communication,
	}
	client.task <- task
	clients[clientType] = client
	return client, nil
}

//往客户端中添加任务
func (c *Client) PutTask(task *downTask) {
	c.task <- task
}

//获取任务
func (c *Client) GetTask() *downTask {
	if len(c.task) <= 0 {
		return nil
	}
	return <-c.task
}

//开始下载
func (c *Client) Download() error {

}

//获取下载的文件信息
func (c *Client) ParseInfo()  {

}
