/*
 + ------------------------------------------------+
 | Author: Zoueature                               |
 + ------------------------------------------------+
 | Email: zoueature@gmail.com                      |
 + ------------------------------------------------+
 | Date: 2020/3/14                                |
 + ------------------------------------------------+
 | Time: 15:13                                     |
 + ------------------------------------------------+
 | Description:                                    |
 + ------------------------------------------------+
*/

package download

import (
	"github.com/zoueature/powerful-download/modules/msg"
	"time"
)

type taskStatus int

type downTask struct {
	taskID       int64             //任务ID用于唯一标示下载任务
	URL          string            //下载URL
	fileName     string            //下载存储文件名
	averageSpeed int               //下载速度， unit bit
	nowSpeed     speed             //实时速度
	breakPointer int               //下载文件偏移量， 用于断点续传
	channel      chan *msg.Message //与客户端的通信通道
	threadNum    int               //下载线程数
	downloader   downloader        //下载器
	status       taskStatus        //任务状态
}

type speed struct {
	startBit     int64
	endBit       int64
	startTime    time.Time
	speedChannel chan<- int64
	bitChannel   <-chan int64
}

const (
	taskStatusWaiting     = 1
	taskStatusDownloading = 2
	taskStatusCancel      = 3
	taskStatusTimeout     = 4
)

func (s *speed) getSpeed() int64 {
	downloadBit := s.endBit - s.startBit
	diffTimestamp := time.Now().Unix() - s.startTime.Unix()
	speed := downloadBit / diffTimestamp
	return speed
}

func (t *downTask) parseURLInfo() error {
	return nil
}

func (t *downTask) listenSignal() {

}

func (t *downTask) reportDownloading() {
	t.channel <- &msg.Message{
		MsgType:    msg.MsgRepportStatus,
		MsgContent: taskStatusDownloading,
	}
}
