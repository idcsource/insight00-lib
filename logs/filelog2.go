// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ]
// Stephen Fire Meditation Qin [ 火志溟 ] -> stephenfmqin@gmail.com
// This source code is governed by GNU LGPL v3 license

package logs

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/idcsource/insight00-lib/base"
)

// 一个符合logs.Logser接口的文件log
type FileLoger struct {
	file *os.File
	lock *sync.RWMutex
}

// 一个符合logs.Logser接口的文件log实例的创建
func NewFileLoger(filename string) (fl *FileLoger, err error) {
	filename = base.LocalFile(filename)
	files, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		err = fmt.Errorf("logs: %v", err)
		return
	}
	fl = &FileLoger{
		file: files,
		lock: new(sync.RWMutex),
	}
	return
}

// 追加一条Log
func (fl *FileLoger) WriteLog(l string) (err error) {
	l = time.Now().Format("2006-01-02 15:04:05.000") + " : " + l + "\n"

	//加写锁
	fl.lock.Lock()
	defer fl.lock.Unlock()

	// 查找文件末尾的偏移量
	theend, err := fl.file.Seek(0, os.SEEK_END)
	if err != nil {
		err = fmt.Errorf("logs: %v", err)
		return
	}
	// 写入一条log
	_, err = fl.file.WriteAt([]byte(l), theend)
	if err != nil {
		err = fmt.Errorf("logs: %v", err)
		return
	}

	return
}

// 读取所有log
func (fl *FileLoger) ReadLogs() (logs []string, err error) {
	//加锁
	fl.lock.RLock()
	defer fl.lock.RUnlock()

	logs = make([]string, 0)

	fi, err := fl.file.Stat()
	if err != nil {
		err = fmt.Errorf("logs: %v", err)
		return
	}
	if fi.Size() > 0 {
		fl.file.Seek(0, 0)

		buf := bufio.NewScanner(fl.file)
		buf.Split(bufio.ScanLines)
		for buf.Scan() {
			fmt.Println(buf.Text())
			logs = append(logs, buf.Text())
		}
	}
	return
}

// 读最后一条log
func (fl *FileLoger) ReadLast() (l string, err error) {
	//加锁
	fl.lock.RLock()
	defer fl.lock.RUnlock()

	fi, err := fl.file.Stat()
	if err != nil {
		err = fmt.Errorf("logs: %v", err)
		return
	}
	if fi.Size() > 0 {
		// 末尾位置
		ret, _ := fl.file.Seek(0, os.SEEK_END)

		seek := int64(0)

		r := make([]byte, 1)
		for i := ret - 2; i >= 0; i-- {
			fl.file.ReadAt(r, i)
			if string(r) == "\n" {
				seek = i
				break
			}
		}

		l_b := make([]byte, ret-seek)
		fl.file.ReadAt(l_b, seek)
		l = string(l_b)
		l = strings.TrimSpace(l)

	}
	return
}

// 关闭清空
func (fl *FileLoger) Close() {
	fl.file.Close()
}

// 写入一条日志
func (fl *FileLoger) Write(p []byte) (n int, err error) {
	n, err = fl.file.Write(p)
	if err != nil {
		err = fmt.Errorf("logs: %v", err)
		return
	}
	return
}
