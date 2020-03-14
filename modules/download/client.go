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
	"github.com/zoueature/powerful-download/modules/msg"
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
	communication map[int64]chan *msg.Message
}

const (
	httpClient    = "http"
	magnetClient  = "magnet"
	torrentClient = "torrent"
)

var (
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
	communication := make(map[int64]chan *msg.Message, config.maxRunTask)
	client := &Client{
		DownConfig:    config,
		task:          taskChan,
		communication: communication,
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
	if c.runTaskNum >= c.maxRunTask {
		c.task <- task
	}
	go task.downloader.download(task)
	c.updateRunTaskNum(1)
	return nil
}

func (c *Client) updateRunTaskNum(step int) {
	c.lock.Lock()
	c.runTaskNum += step
	c.lock.Unlock()
}

func (c *Client) InitTask(url string) (*downTask, error) {
	downloader := parseToDownloader(url)
	if downloader == nil {
		return nil, errors.New("generate downloader error ")
	}
	task := &downTask{
		taskID: time.Now().UnixNano(),
		URL: url,
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
	channel := make(chan *msg.Message)
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
