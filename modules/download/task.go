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

import "time"

type downTask struct {
	taskID       int64
	URL          string
	fileName     string
	averageSpeed int //unit bit
	nowSpeed     speed
	breakPointer int //下载文件偏移量， 用于断点续传
	channel      chan interface{}
}
type speed struct {
	startBit     int64
	endBit       int64
	startTime    time.Time
	speedChannel chan<- int64
	bitChannel   <-chan int64
}

func (s *speed) getSpeed() int64 {
	downloadBit := s.endBit - s.startBit
	diffTimestamp := time.Now().Unix() - s.startTime.Unix()
	speed := downloadBit / diffTimestamp
	return speed
}



