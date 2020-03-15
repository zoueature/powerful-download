/*+-----------------------------+
 *| Author: Zoueature           |
 *+-----------------------------+
 *| Email: zoueature@gmail.com  |
 *+-----------------------------+
 */
package model

const (
	TaskStatusWaitingDownloaded = 1
	TaskStatusDone              = 2
	TaskStatusDelete            = 3

	TaskTypeHTTP    = 1
	TaskTypeMagnet  = 2
	TaskTypeTorrent = 3
	TaskTypeED2K    = 4
)

type Task struct {
	ID           int64  `gorm:"column:id"`
	URL          string `gorm:"column:url"`
	FileName     string `gorm:"column:fileName"`
	Type         int    `gorm:"column:type"`
	BreakPointer int64  `gorm:"column:break_pointer"`
	Status       int    `form:"column:status"`
	CreateTime   int    `gorm:"column:create_time"`
	FinishTime   int    `gorm:"column:finish_time"`
}

func (t *Task) TableName() string {
	return taskDbTable
}

// 获取所有任务
func (conn *connection) GetTasks(status int) []*Task {
	tasks := make([]*Task, 0)
	var where Task
	if status != 0 {
		where.Status = status
	}
	result := conn.db.Where(where).Find(&tasks)
	if result.Error != nil {
		return nil
	}
	return tasks
}

//创建任务
func (conn *connection) CreateTask(ID int64, Type int, URL, fileName string) error {
	task := Task{
		ID:       ID,
		URL:      URL,
		FileName: fileName,
		Type:     Type,
		Status:   TaskStatusWaitingDownloaded,
	}
	result := conn.db.Create(task)
	return result.Error
}
