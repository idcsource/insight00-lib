// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ]
// Stephen Fire Meditation Qin [ 火志溟 ] -> stephenfmqin@gmail.com
// This source code is governed by GNU LGPL v3 license

package logs

import (
	"log"
	"os"
)

// 接口实现
type Logser interface {
	WriteLog(log string) (err error)      // 写一条log
	ReadLogs() (logs []string, err error) // 读全部log
	ReadLast() (l string, err error)      // 读最后一条
	Close()                               //关闭
}

// 日志
type Logs struct {
	runtimelog *RuntimeLog
	logs       *log.Logger
}

// Runtime log, in the memory, if software stop, the log will disappear.
type RuntimeLog struct {
	logs   []string
	maxnum int
}

// 文件型日志
type FileLog struct {
	file *os.File
}
