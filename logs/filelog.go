// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package logs

import (
	"fmt"
	"os"
	"time"

	"github.com/idcsource/insight00-lib/base"
)

// 创建一个新的文件型日志
func NewFileLog(filename string) (fl *FileLog, err error) {
	filename = base.LocalFile(filename)
	files, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		err = fmt.Errorf("logs: %v", err)
		return
	}
	fl = &FileLog{
		file: files,
	}
	return
}

// 写入一条标准格式的日志
func (fl *FileLog) WriteLog(log string) (err error) {
	log = time.Now().Format("2006-01-02 15:04:05.000") + " : " + log + "\n"
	_, err = fl.file.Write([]byte(log))
	if err != nil {
		err = fmt.Errorf("logs: %v", err)
	}
	return
}

// 写入一条日志
func (fl *FileLog) Write(p []byte) (n int, err error) {
	n, err = fl.file.Write(p)
	if err != nil {
		err = fmt.Errorf("logs: %v", err)
		return
	}
	return
}
