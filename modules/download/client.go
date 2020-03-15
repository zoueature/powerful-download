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
	"fmt"
	"github.com/zoueature/powerful-download/modules/msg"
	"log"
	"strings"
	"sync"
	"time"
)

type downloader interface {
	download(task *downTask) error
	cancel(task *downTask) error
	parseURLInfo(url string) (*downloadInfo, error)
}

type urlChecker func(string) bool
type downloadInstance func() downloader
type downloadMaker struct {
	instance downloader
	check    urlChecker
	make     downloadInstance
}

type Client struct {
	DownConfig
	task          chan *downTask
	lock          sync.RWMutex
	runTaskNum    int
	communication map[int64]chan *msg.TaskMsg
	taskChan      chan *msg.TaskMsg
	operatorChan  chan *msg.OperatorMsg
}

const (
	httpClient    = "http"
	magnetClient  = "magnet"
	torrentClient = "torrent"
)

var (
	//支持的url格式
	urlSupport = map[string]downloadMaker{
		httpClient:    {nil, isHttp, newHTTP},
		magnetClient:  {nil, isMagnet, newMagnet},
		torrentClient: {nil, isTorrent, newTorrent},
	}
)

//实例化一个下载客户端
func NewClient(config DownConfig) (*Client, error) {
	if config.maxSpeed == 0 || config.maxRunTask == 0 || config.maxTask == 0 {
		return nil, errors.New("error config ")
	}
	//任务队列
	taskChan := make(chan *downTask, config.maxTask)
	//通信通道
	communication := make(map[int64]chan *msg.TaskMsg)
	client := &Client{
		DownConfig:    config,
		task:          taskChan,
		communication: communication,
		taskChan:      make(chan *msg.TaskMsg),
		operatorChan:  make(chan *msg.OperatorMsg),
	}
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
func (c *Client) Download(task *downTask) error {
	err := task.InsertToDB()
	if err != nil {
		return err
	}
	if c.runTaskNum >= c.maxRunTask {
		c.task <- task
	}
	task.downloader.download(task)
	c.updateRunTaskNum(1)
	return nil
}

//更新正在下载数量
func (c *Client) updateRunTaskNum(step int) {
	c.lock.Lock()
	c.runTaskNum += step
	c.lock.Unlock()
}

//初始化下载任务
func (c *Client) InitTask(url string) (*downTask, error) {
	downloader := parseToDownloader(url)
	if downloader == nil {
		return nil, errors.New("generate downloader error ")
	}
	task := &downTask{
		taskID: time.Now().UnixNano(),
		URL:    url,
	}
	info, err := downloader.parseURLInfo(url)
	if err != nil {
		return nil, err
	}
	task.fileName = info.fileName
	if info.supportPartDownload {
		task.threadNum = c.threadNum
	} else {
		task.threadNum = 1
	}
	task.downloader = downloader
	channel := make(chan *msg.TaskMsg)
	task.channel = channel
	c.communication[task.taskID] = channel
	return task, nil
}

func parseToDownloader(url string) downloader {
	url = strings.Trim(url, " ")
	for _, maker := range urlSupport {
		if maker.check(url) {
			if maker.instance == nil {
				maker.instance = maker.make()
			}
			return maker.instance
		}
	}
	return nil
}

func (c *Client) Run() {
	go func() {
		for {
			fmt.Println("listen")
			select {
			case m := <-c.taskChan:
				//task任务消息
				err := c.handleTaskMsg(m)
				if err != nil {
					log.Println(err.Error())
				}
			case m := <-c.operatorChan:
				//处理操作者消息
				fmt.Println(m)
				err := c.handleOperateMsg(m)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
	}()
}

//处理收到的task任务的消息
func (c *Client) handleOperateMsg(m *msg.OperatorMsg) error {
	switch m.MsgType {
	case msg.OperateMsgNewDownload:
		url, ok := m.MsgContent.(string)
		if !ok || url == "" {
			return errors.New("get url info error ")
		}
		fmt.Println(url)
		task, err := c.InitTask(url)
		if err != nil {
			return err
		}
		err = c.Download(task)
		if err != nil {
			return err
		}
	}
	return nil
}

//处理用户发出的操作信号
func (c *Client) handleTaskMsg(m *msg.TaskMsg) error {
	return nil
}

//发送信号给task任务
func (c *Client) SendTaskMsg(taskID int64, m *msg.TaskMsg) error {
	return nil
}

//发送消息给用户
func (c *Client) SendOperateMsg(m *msg.OperatorMsg) error {
	return nil
}
