// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package logs

type LogsType uint8

const (
	LOGS_TYPE_IS_RUNTIME_LOG LogsType = iota
	LOGS_TYPE_IS_FILE_LOG
)

const (
	FILE_LOGER_MAX_LOG    uint64 = 1000
	FILE_LOGER_WRITE_CHAN int    = 10
)
