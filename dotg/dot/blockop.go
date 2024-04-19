// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ]
// Stephen Fire Meditation Qin [ 火志溟 ] -> stephenfmqin@gmail.com
// This source code is governed by GNU LGPL v3 license

package dot

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/idcsource/insight00-lib/base"
	"github.com/idcsource/insight00-lib/iendecode"
)

// 对某个block进行操作
type BlockOp struct {
	path           string              // block的位置，这两个都要与InitBlock一致
	version        uint8               // block的版本
	deep           uint8               // block的路径深度，这两个都要与InitBlock一致
	dots_lock      map[string]*DotLock // 正在操作的dot都会加上相应的锁，map的key为dot的id
	dots_lock_lock *sync.RWMutex       // 避免操作上面的dot锁时有抢占，在对上面的锁修改时也要现锁定
}

// dot的操作锁
type DotLock struct {
	LockTime time.Time
	LockType _BlockDotLockType // 锁的类型
	Lock     *sync.RWMutex
}

// 新建一个dot
func (bop *BlockOp) NewDot(id string, data []byte) (fpath string, fname string, err error) {
	fname, fpath, err = bop.findFilePath(id)
	if err != nil {
		err = fmt.Errorf("%v", err)
		return
	}

	// 加锁
	bop.dots_lock_lock.Lock()
	if _, have := bop.dots_lock[id]; have != true {
		bop.dots_lock[id] = &DotLock{
			LockTime: time.Now(),
			LockType: BLOCK_DOT_LOCK_TYPE_NOTHING,
			Lock:     new(sync.RWMutex),
		}
	}
	// 如果没有锁就加内部锁，如果是外部锁，就不管了
	if bop.dots_lock[id].LockType == BLOCK_DOT_LOCK_TYPE_NOTHING {
		bop.dots_lock[id].LockTime = time.Now()
		bop.dots_lock[id].LockType = BLOCK_DOT_LOCK_TYPE_INSIDE
		bop.dots_lock[id].Lock.Lock()
		defer func() {
			bop.dots_lock_lock.Lock()
			bop.dots_lock[id].Lock.Unlock()
			bop.dots_lock[id].LockType = BLOCK_DOT_LOCK_TYPE_NOTHING
			bop.dots_lock_lock.Unlock()
		}()
	}
	bop.dots_lock_lock.Unlock()

	// 确认文件
	ishave_body := base.FileExist(fpath + fname + "_body")
	ishave_context_index := base.FileExist(fpath + fname + "_context_index")
	ishave_context_del_index := base.FileExist(fpath + fname + "_context_del_index")
	if ishave_body == true || ishave_context_index == true || ishave_context_del_index == true {
		err = fmt.Errorf("Dot Block: The dot id \"%v\" already have.", id)
		return
	}
	// 准备基本头部数据
	dotversion_b := iendecode.Uint8ToBytes(DOT_NOW_DEFAULT_VERSION) // dot程序版本
	opversion_b := iendecode.Uint64ToBytes(0)                       // 操作版本

	//打开写入body文件
	dot_body_f, err := os.OpenFile(fpath+fname+"_body", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	defer dot_body_f.Close()
	//打开写入context索引文件
	dot_context_index_f, err := os.OpenFile(fpath+fname+"_context_index", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	defer dot_context_index_f.Close()
	//打开写入context删除索引文件
	dot_context_del_index_f, err := os.OpenFile(fpath+fname+"_context_del_index", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	defer dot_context_del_index_f.Close()

	// 开始写body
	_, err = dot_body_f.Write(dotversion_b) // 应用版本
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	_, err = dot_body_f.Write(bop.idToByte255(id)) // ID
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	_, err = dot_body_f.Write(opversion_b) // 操作版本
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	_, err = dot_body_f.Write(data) // 数据
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	// 开始写context索引文件，现在还是一个空文件
	_, err = dot_context_index_f.Write(dotversion_b) // 应用版本
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	_, err = dot_context_index_f.Write(opversion_b) // 操作版本
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	// 开始写context删除索引，现在还是一个空文件
	_, err = dot_context_del_index_f.Write(dotversion_b) // 应用版本
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	_, err = dot_context_del_index_f.Write(opversion_b) // 操作版本
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	return
}

// 显示当前的全部dot锁状态
func (bop *BlockOp) DisplayDotLock() (dots_lock map[string]*DotLock) {
	return bop.dots_lock
}

// 外部加锁
func (bop *BlockOp) OutLock(id string) (err error) {
	return

}

// 外部加读锁
func (bop *BlockOp) OutRLock(id string) (err error) {
	return

}

// 外部解锁
func (bop *BlockOp) OutUnlock(id string) (err error) {
	return

}

// 外部读解锁
func (bop *BlockOp) OutRUnlock(id string) (err error) {
	return

}

// 返回要操作的dot的文件名和路径，同时会检查ID的长度
func (bop *BlockOp) findFilePath(id string) (fname string, fpath string, err error) {
	if len([]byte(id)) > DOT_ID_MAX_LENGTH_V2 {
		err = fmt.Errorf("Dot Block: The dot id length must less than %v: \"%v\"", DOT_ID_MAX_LENGTH_V2, id)
		return
	}
	fname = base.GetSha1Sum(id)
	fpath = bop.path
	for i := 0; i < int(bop.deep); i++ {
		fpath = fpath + string(fname[i]) + "/"
	}

	return
}

func (bop *BlockOp) idToByte255(id string) (b []byte) {
	id_b := []byte(id)

	b = make([]byte, DOT_ID_MAX_LENGTH_V2)
	for i := 0; i < len(id_b); i++ {
		b[i] = id[i]
	}
	return
}

func (bop *BlockOp) byte255ToId(b []byte) (id string) {
	var id_b []byte
	for j := 0; j < DOT_ID_MAX_LENGTH_V2; j++ {
		if b[j] != 0 {
			id_b = append(id_b, b[j])
		}
	}
	id = string(id_b)
	return
}

// 从文件里读取多少以后的全部数据
func (bop *BlockOp) readAfter(m int64, fname string) (b []byte, len int64, err error) {
	f, err := os.Open(fname)
	if err != nil {
		return
	}
	defer f.Close()

	var size int64
	if info, err := f.Stat(); err == nil {
		size = info.Size()
	}
	len = size - m
	if len <= 0 {
		b = make([]byte, 0)
		return
	}
	b = make([]byte, len)
	_, err = f.ReadAt(b, m)
	return
}

// 从文件里读取多少以后的全部数据
func (bop *BlockOp) readAfterWithFile(m int64, f *os.File) (b []byte, len int64, err error) {
	var size int64
	if info, err := f.Stat(); err == nil {
		size = info.Size()
	}
	len = size - m
	if len <= 0 {
		b = make([]byte, 0)
		return
	}
	b = make([]byte, len)
	_, err = f.ReadAt(b, m)
	return
}
