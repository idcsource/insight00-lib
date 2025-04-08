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
	filename     string
	max          uint64 // 单文件最大日志条目数
	now          uint64
	file         *os.File
	lock         *sync.RWMutex
	writechannel chan string
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
		filename: filename,
		max:      FILE_LOGER_MAX_LOG,
		now:      0,
		file:     files,
		lock:     new(sync.RWMutex),
	}
	fl.writechannel = make(chan string, FILE_LOGER_WRITE_CHAN)

	// 查看当前日志文件里面的条目数
	fi, err := files.Stat()
	if err != nil {
		err = fmt.Errorf("logs: %v", err)
		return
	}
	if fi.Size() > 0 {
		files.Seek(0, 0)

		buf := bufio.NewScanner(files)
		buf.Split(bufio.ScanLines)

		for buf.Scan() {
			fl.now++
		}
		if err = buf.Err(); err != nil {
			err = fmt.Errorf("logs: %v", err)
			return
		}
	}

	// 如果文件大了，就新建
	if fl.now >= fl.max {
		err = fl.newLogFile()
		if err != nil {
			err = fmt.Errorf("logs: %v", err)
			return
		}
	}

	go fl.goToLog() // 非堵塞的写日志

	return
}

// 写log的go
func (fl *FileLoger) goToLog() {
	for {
		select {
		case thelog := <-fl.writechannel:
			fl.toWriteLog(thelog)
		default:

		}

	}
}

// 实际的写日志函数
func (fl *FileLoger) toWriteLog(thelog string) (err error) {
	//加写锁
	fl.lock.Lock()
	defer fl.lock.Unlock()

	// 如果文件大了，就新建
	if fl.now >= fl.max {
		err = fl.newLogFile()
		if err != nil {
			err = fmt.Errorf("logs: %v", err)
			return
		}
	}

	fl.now++

	// 查找文件末尾的偏移量
	theend, err := fl.file.Seek(0, os.SEEK_END)
	if err != nil {
		err = fmt.Errorf("logs: %v", err)
		return
	}
	// 写入一条log
	_, err = fl.file.WriteAt([]byte(thelog), theend)
	if err != nil {
		err = fmt.Errorf("logs: %v", err)
		return
	}
	return
}

// 新建一个log文件
func (fl *FileLoger) newLogFile() (err error) {
	//关闭文件
	err = fl.file.Close()
	if err != nil {
		err = fmt.Errorf("logs: %v", err)
		return
	}
	//重命名文件
	rename := fl.filename + "_" + time.Now().Format("2006-01-02_15:04:05")
	err = os.Rename(fl.filename, rename)
	if err != nil {
		err = fmt.Errorf("logs: %v", err)
		return
	}
	//创建新文件
	files, err := os.OpenFile(fl.filename, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		err = fmt.Errorf("logs: %v", err)
		return
	}
	fl.file = files
	fl.now = 0
	return
}

// 设置每个日志文件最多存储的日志条数，默认是1000
func (fl *FileLoger) SetCountPerFile(max_count uint64) {
	fl.max = max_count
}

// 追加一条Log
func (fl *FileLoger) WriteLog(l string) (err error) {
	l = time.Now().Format("2006-01-02 15:04:05.000") + " : " + l + "\n"

	fl.writechannel <- l

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
