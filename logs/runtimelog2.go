// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ]
// Stephen Fire Meditation Qin [ 火志溟 ] -> stephenfmqin@gmail.com
// This source code is governed by GNU LGPL v3 license

package logs

import (
	"sync"
	"time"
)

// 一个符合logs.Logser接口的runtimelog
type RunLoger struct {
	max  int64
	now  int64 // 这是下次要写的位置
	logs []string
	lock *sync.RWMutex
}

// 一个符合logs.Logser接口的runtimelog实例的创建，max为最大记录log的条数
func NewRunLoger(max int64) (rl *RunLoger, err error) {
	rl = &RunLoger{
		max:  max,
		now:  0,
		logs: make([]string, max),
		lock: new(sync.RWMutex),
	}
	return
}

// 追加一条Log
func (rl *RunLoger) WriteLog(l string) (err error) {
	l = time.Now().Format("2006-01-02 15:04:05.000") + " : " + l

	//加写锁
	rl.lock.Lock()
	defer rl.lock.Unlock()

	// 写入一条log
	rl.logs[rl.now] = l
	rl.now++
	if rl.now == rl.max {
		rl.now = 0
	}

	return
}

// 读取所有log
func (rl *RunLoger) ReadLogs() (logs []string, err error) {

	//加锁
	rl.lock.RLock()
	defer rl.lock.RUnlock()

	logs = make([]string, 0)

	n := rl.now
	if n == rl.max {
		n = 0
	}
	for {
		if rl.logs[n] != "" {
			logs = append(logs, rl.logs[n])
		}
		n++
		if n == rl.now {
			break
		}
		if n == rl.max {
			n = 0
		}
	}

	return
}

// 读最后一条log
func (rl *RunLoger) ReadLast() (l string, err error) {

	//加锁
	rl.lock.RLock()
	defer rl.lock.RUnlock()

	l = rl.logs[rl.now-1]

	return
}

// 关闭清空
func (rl *RunLoger) Close() {
	rl.now = 0
	rl.logs = make([]string, rl.max)
}

// 写入一条运行时日志
func (rl *RunLoger) Write(p []byte) (n int, err error) {

	l := string(p)

	//加写锁
	rl.lock.Lock()
	defer rl.lock.Unlock()

	// 写入一条log
	rl.logs[rl.now] = l
	rl.now++
	if rl.now == rl.max {
		rl.now = 0
	}

	n = len(string(p))
	return
}
