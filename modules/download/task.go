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
	"errors"
	"fmt"
	"github.com/zoueature/powerful-download/model"
	"github.com/zoueature/powerful-download/modules/msg"
	"log"
	"time"
)

type taskStatus int

type downTask struct {
	taskID                 int64                       //任务ID用于唯一标示下载任务
	URL                    string                      //下载URL
	fileName               string                      //下载存储文件名
	averageSpeed           int                         //下载速度， unit bit
	nowSpeed               speed                       //实时速度
	breakPointer           int                         //下载文件偏移量， 用于断点续传
	channel                chan *msg.TaskMsg           //与客户端的通信通道
	threadNum              int                         //下载线程数
	downloader             downloader                  //下载器
	status                 taskStatus                  //任务状态
	commander              chan *msg.TaskMsg           //任务下载程序与监听程序的通信通道
	communication          map[int64]chan *msg.TaskMsg //task线程与下载线程的通信
	threadNumInDownloading int                         //正在下载的线程数
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
	go func() {
		for {
			fmt.Println("task is listening signal")
			select {
			case m := <-t.channel:
				//客户端消息
				err := t.handleClientMsg(m)
				if err != nil {
					log.Println(err.Error())
				}
			case m := <-t.commander:
				//todo 下载线程消息
				err := t.handleDownloadThreadMsg(m)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
	}()
}

//处理客户端发来的消息
func (t *downTask) handleClientMsg(m *msg.TaskMsg) error {
	switch m.MsgType {
	case msg.TaskMsgStop:
		if len(t.communication) <= 0 {
			//无下载线程直接返回
			return nil
		}
		//通知下载线程停止下载
		stopMsg := &msg.TaskMsg{}
		stopMsg.MsgType = msg.TaskMsgStop
		for _, ch := range t.communication {
			ch <- stopMsg
		}
	}
	return nil
}

//处理下载线程发来的消息
func (t *downTask) handleDownloadThreadMsg(m *msg.TaskMsg) error {
	switch m.MsgType {

	}
	return nil
}

func (t *downTask) InsertToDB() error {
	conn := model.GetModel()
	var ty int
	switch t.downloader.(type) {
	case *HttpDownloader:
		ty = model.TaskTypeHTTP
	case *MagnetDownloader:
		ty = model.TaskTypeMagnet
	case *TorrentDownloader:
		ty = model.TaskTypeTorrent
	}
	if ty == 0 {
		return errors.New("不支持的下载类型")
	}
	err := conn.CreateTask(t.taskID, ty, t.URL, t.fileName)
	return err
}
