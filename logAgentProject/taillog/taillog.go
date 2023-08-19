package taillog

import (
	"context"
	"fmt"
	"time"

	"github.com/Cc-Love-Study/LogAgent/kafka"
	"github.com/hpcloud/tail"
)

// 一个日志搜集的任务
type TailTask struct {
	Path       string
	Topic      string
	TailObj    *tail.Tail
	Ctx        context.Context
	CancelFunc context.CancelFunc
}

func NewTailTask(path string, topic string) (*TailTask, error) {
	ctx, cancel := context.WithCancel(context.Background())
	tailTask := &TailTask{
		Path:       path,
		Topic:      topic,
		Ctx:        ctx,
		CancelFunc: cancel,
	}

	err := tailTask.tailInit()
	if err != nil {
		return nil, err
	}
	return tailTask, nil
}

func (t *TailTask) tailInit() (err error) {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	t.TailObj, err = tail.TailFile(t.Path, config)
	if err != nil {
		return
	}
	go t.run()
	return nil
}

// 将数据放入到通道缓冲区
// 不是直接使用kafka进行发送
func (t *TailTask) run() {
	for {
		select {
		case line := <-t.TailObj.Lines:
			kafka.SendToChan(t.Topic, line.Text)
		case <-t.Ctx.Done():
			fmt.Printf("task %s_%s 退出", t.Path, t.Topic)
			return
		default:
			time.Sleep(time.Second)
		}
	}

}
