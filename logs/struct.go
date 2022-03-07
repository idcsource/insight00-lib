// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package logs

import (
	"log"
	"os"
)

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
