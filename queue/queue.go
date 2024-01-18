package queue

import (
	"time"

	"github.com/lib/pq"
)

// Queue 队列
type Queue struct {
	listener *pq.Listener
}

// NewQueue 构造函数
// channel: 创建postgresql函数时的通道名称
// dsn: postgresql连接字符串, 例: dbname=dingtalk_server user=postgres password=postgres host=127.0.0.1 sslmode=disable
func NewQueue(dsn, channel string) (*Queue, error) {
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			panic(err)
		}
	}

	listener := pq.NewListener(dsn, 10*time.Second, time.Minute, reportProblem)
	if err := listener.Listen(channel); err != nil {
		return nil, err
	}
	return &Queue{listener: listener}, nil
}

// WaitForNotification 等待通知
func (q *Queue) WaitForNotification() <-chan string {
	notice := make(chan string)
	go func() {
		for {
			select {
			case n := <-q.listener.Notify:
				notice <- n.Extra
			case <-time.After(90 * time.Second):
				go func() {
					q.listener.Ping()
				}()
			}
		}
	}()
	return notice
}
