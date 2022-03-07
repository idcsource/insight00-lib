// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package logs

import (
	"fmt"
	"log"
)

// 创建日志
func NewLogs() (logs *Logs) {
	logs = &Logs{
		runtimelog: make(map[string]*RuntimeLog),
		logs:       make(map[string]*log.Logger),
	}
	return
}

// 设置一个运行时日志
func (logs *Logs) SetRuntimeLog(name string, prefix string, maxnum int) (err error) {
	rtl := NewRuntimeLog(maxnum)
	prefix = prefix + " "
	runlogs := log.New(rtl, prefix, log.Ldate|log.Ltime)
	logs.runtimelog[name] = rtl
	logs.logs[name] = runlogs
	return
}

// 设置一个文件日志
func (logs *Logs) SetFileLog(name string, prefix string, filename string) (err error) {
	fl, err := NewFileLog(filename)
	if err != nil {
		return
	}
	prefix = prefix + " "
	flogs := log.New(fl, prefix, log.Ldate|log.Ltime)
	logs.logs[name] = flogs
	return
}

// 写入日志
func (logs *Logs) PrintLog(name string, s ...interface{}) (err error) {
	_, have := logs.logs[name]
	if have == false {
		err = fmt.Errorf("logs: Thers's no log name \"%v\"", name)
		return
	}
	logs.logs[name].Print(s)
	return
}

// 输入运行时日志
func (logs *Logs) RuntimeOutput(name string) (output []string, err error) {
	_, have := logs.runtimelog[name]
	if have == false {
		err = fmt.Errorf("logs: Thers's no log name \"%v\"", name)
		return
	}
	output = logs.runtimelog[name].Output()
	return
}
