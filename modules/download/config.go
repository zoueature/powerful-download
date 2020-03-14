/*
 + ------------------------------------------------+
 | Author: Zoueature                               |
 + ------------------------------------------------+
 | Email: zoueature@gmail.com                      |
 + ------------------------------------------------+
 | Date: 2020/3/14                                |
 + ------------------------------------------------+
 | Time: 14:17                                     |
 + ------------------------------------------------+
 | Description:                                    |
 + ------------------------------------------------+
*/

package download

const (
	defaultThreadNumConfig    = 10
	defaultDownloadPathConfig = "./powerful-download/"
	defaultTimeoutConfig      = 10
	defaultMaxSpeedConfig     = 1024 * 1024 * 1034 //1M
	defaultMaxRunTask         = 10
	defaultMaxTask            = 100
)

//下载配置
type DownConfig struct {
	threadNum      int    //线程数
	toPath         string //下载路劲
	maxSpeed       int    //最高下载速度
	connectTimeout int    //连接超时时间
	maxRunTask     int    //并行下载的最大任务数
	maxTask        int    //队列中的最高任务数
}

func NewConfig(threadNum int, toPath string, maxSpeed, connectTimeout int) DownConfig {
	if threadNum == 0 {
		threadNum = defaultThreadNumConfig
	}
	if toPath == "" {
		toPath = defaultDownloadPathConfig
	}
	if maxSpeed == 0 {
		maxSpeed = defaultMaxSpeedConfig
	}
	if connectTimeout == 0 {
		connectTimeout = defaultTimeoutConfig
	}
	return DownConfig{
		threadNum:      threadNum,
		toPath:         toPath,
		maxSpeed:       maxSpeed,
		connectTimeout: connectTimeout,
		maxTask:        defaultMaxTask,
		maxRunTask:     defaultMaxRunTask,
	}
}
