// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package logs

import (
	"log"
	"os"
)

type Logs struct {
	logstype LogsType
	log      *log.Logger
	logstart bool
}

// Runtime log, in the memory, if software stop, the log will disappear.
type RuntimeLogContent struct {
	logs   []string
	maxnum int
}

type FileLogContent struct {
	file *os.File
}
