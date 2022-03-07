// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package logs

import (
	"log"
)

// 创建日志
func NewLogs() (logs *Logs) {
	logs = &Logs{}
	return
}

// 设置一个运行时日志
func (logs *Logs) SetRuntimeLog(prefix string, maxnum int) (err error) {
	rtl := NewRuntimeLog(maxnum)
	prefix = prefix + " "
	runlogs := log.New(rtl, prefix, log.Ldate|log.Ltime)
	logs.runtimelog = rtl
	logs.logs = runlogs
	return
}

// 设置一个文件日志
func (logs *Logs) SetFileLog(prefix string, filename string) (err error) {
	fl, err := NewFileLog(filename)
	if err != nil {
		return
	}
	prefix = prefix + " "
	flogs := log.New(fl, prefix, log.Ldate|log.Ltime)
	logs.logs = flogs
	return
}

// 写入日志
func (logs *Logs) PrintLog(s ...interface{}) (err error) {
	logs.logs.Print(s)
	return
}

// 输入运行时日志
func (logs *Logs) RuntimeOutput() (output []string, err error) {
	output = logs.runtimelog.Output()
	return
}
