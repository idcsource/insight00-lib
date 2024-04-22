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

// 修改dot中的数据
func (bop *BlockOp) UpdateDotData(dotid string, data []byte) (err error) {
	fname, fpath, err := bop.findFilePath(dotid)
	if err != nil {
		err = fmt.Errorf("%v", err)
		return
	}

	//加锁
	bop.dots_lock_lock.Lock()
	if _, have := bop.dots_lock[dotid]; have != true {
		bop.dots_lock[dotid] = &DotLock{
			LockTime: time.Now(),
			LockType: BLOCK_DOT_LOCK_TYPE_NOTHING,
			Lock:     new(sync.RWMutex),
		}
	}
	// 如果没有锁就加内部锁，如果是外部锁，就不管了
	if bop.dots_lock[dotid].LockType == BLOCK_DOT_LOCK_TYPE_NOTHING {
		bop.dots_lock[dotid].LockTime = time.Now()
		bop.dots_lock[dotid].LockType = BLOCK_DOT_LOCK_TYPE_INSIDE
		bop.dots_lock[dotid].Lock.Lock()
		defer func() {
			bop.dots_lock_lock.Lock()
			bop.dots_lock[dotid].Lock.Unlock()
			bop.dots_lock[dotid].LockType = BLOCK_DOT_LOCK_TYPE_NOTHING
			bop.dots_lock_lock.Unlock()
		}()
	}
	bop.dots_lock_lock.Unlock()

	// 构建文件名
	fname_data := fname + "_body"
	ishave := base.FileExist(fpath + fname_data)
	// 如果不存在就返回错误
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the dot \"%v\".", dotid)
		return
	}

	// 打开数据文件写入
	bop_data_f, err := os.OpenFile(fpath+fname_data, os.O_RDWR, 0600)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	defer bop_data_f.Close()

	// 获取操作版本，并且+1
	opversion_b := make([]byte, 8)
	read_n, err := bop_data_f.ReadAt(opversion_b, 1+DOT_ID_MAX_LENGTH_V2)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	if read_n != 8 {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	opversion := iendecode.BytesToUint64(opversion_b)
	opversion++
	opversion_b = iendecode.Uint64ToBytes(opversion)

	// 准备写入更新的数据
	// 扔掉之前的数据部分
	err = bop_data_f.Truncate(1 + DOT_ID_MAX_LENGTH_V2 + 8)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	// 写新的操作版本
	_, err = bop_data_f.WriteAt(opversion_b, 1+DOT_ID_MAX_LENGTH_V2)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	// 写新的数据
	_, err = bop_data_f.WriteAt(data, 1+DOT_ID_MAX_LENGTH_V2+8)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	return
}

// 读取dot中的数据
func (bop *BlockOp) ReadDotData(dotid string) (data []byte, len int64, err error) {
	fname, fpath, err := bop.findFilePath(dotid)
	if err != nil {
		err = fmt.Errorf("%v", err)
		return
	}

	// 加锁，读锁
	bop.dots_lock_lock.Lock()
	if _, have := bop.dots_lock[dotid]; have != true {
		bop.dots_lock[dotid] = &DotLock{
			LockTime: time.Now(),
			LockType: BLOCK_DOT_LOCK_TYPE_NOTHING,
			Lock:     new(sync.RWMutex),
		}
	}
	// 如果没有锁就加内部锁，如果是外部锁，就不管了
	if bop.dots_lock[dotid].LockType == BLOCK_DOT_LOCK_TYPE_NOTHING {
		bop.dots_lock[dotid].LockTime = time.Now()
		bop.dots_lock[dotid].LockType = BLOCK_DOT_LOCK_TYPE_INSIDE
		bop.dots_lock[dotid].Lock.RLock()
		defer func() {
			bop.dots_lock_lock.Lock()
			bop.dots_lock[dotid].Lock.RUnlock()
			bop.dots_lock[dotid].LockType = BLOCK_DOT_LOCK_TYPE_NOTHING
			bop.dots_lock_lock.Unlock()
		}()
	}
	bop.dots_lock_lock.Unlock()

	// 打开文件
	fname_data := fname + "_body"
	ishave := base.FileExist(fpath + fname_data)
	// 如果不存在就返回错误
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the dot \"%v\".", dotid)
		return
	}
	f, err := os.OpenFile(fpath+fname_data, os.O_RDONLY, 0600)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	data, len, err = bop.readAfterWithFile(1+DOT_ID_MAX_LENGTH_V2+8, f)

	return
}

// 删除一个dot
func (bop *BlockOp) DelDot(dotid string) (err error) {
	return bop.DropDot(dotid)
}
func (bop *BlockOp) DropDot(dotid string) (err error) {
	fname, fpath, err := bop.findFilePath(dotid)
	if err != nil {
		err = fmt.Errorf("%v", err)
		return
	}

	//加锁
	bop.dots_lock_lock.Lock()
	if _, have := bop.dots_lock[dotid]; have != true {
		bop.dots_lock[dotid] = &DotLock{
			LockTime: time.Now(),
			LockType: BLOCK_DOT_LOCK_TYPE_NOTHING,
			Lock:     new(sync.RWMutex),
		}
	}
	// 如果没有锁就加内部锁，如果是外部锁，就不管了
	if bop.dots_lock[dotid].LockType == BLOCK_DOT_LOCK_TYPE_NOTHING {
		bop.dots_lock[dotid].LockTime = time.Now()
		bop.dots_lock[dotid].LockType = BLOCK_DOT_LOCK_TYPE_INSIDE
		bop.dots_lock[dotid].Lock.Lock()
		defer func() {
			bop.dots_lock_lock.Lock()
			bop.dots_lock[dotid].Lock.Unlock()
			bop.dots_lock[dotid].LockType = BLOCK_DOT_LOCK_TYPE_NOTHING
			delete(bop.dots_lock, dotid) // 既然删除了dot，那就不用保留这个锁了
			bop.dots_lock_lock.Unlock()
		}()
	}
	bop.dots_lock_lock.Unlock()

	// 看存不存在
	ishave := base.FileExist(fpath + fname + "_body")
	// 如果不存在这个文件，正好就什么都不处理
	if ishave != true {
		return
	}

	// 读context索引
	context_b, context_b_l, err := bop.readAfter(1+8, fpath+fname+"_context_index")
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	// 如果没有就不用管了，如果有就要把context的列表都读出来
	if context_b_l != 0 {
		context_index := bop.readContextIndex(context_b)
		// 开始删文件了
		for cindex := range context_index {
			the_c_name := fpath + fname + "_context_" + base.GetSha1Sum(context_index[cindex].ContextName)
			os.Remove(the_c_name)
			the_c_d_name := fpath + fname + "_context_" + base.GetSha1Sum(context_index[cindex].ContextName) + "_del_index"
			os.Remove(the_c_d_name)
		}
	}
	// 开始删文件了
	os.Remove(fpath + fname + "_context_index")
	os.Remove(fpath + fname + "_context_del_index")
	os.Remove(fpath + fname + "_body")

	return
}

// 增加一个context

// 删除一个context

// 读取所有context名称

// 修改一个context的up信息（名称+数据）

// 修改一个context的up信息（只名称）

// 修改一个context的up信息（只数据）

// 增加一个context的down信息（名称+数据）

// 修改一个context的down信息（只数据）

// 删除一个context的down信息

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

// 完整读取context的索引
func (bop *BlockOp) readContextIndex(b []byte) (index []ContextIndex) {
	index = make([]ContextIndex, 0)
	var i int64
	b_len := int64(len(b))
	for i = 0; i < b_len; i++ {
		// 状态
		status_b := b[i:1]
		status_uint := iendecode.BytesToUint8(status_b)
		status := _DotContextIndex_Status(status_uint)
		name_b := b[i+1 : DOT_ID_MAX_LENGTH_V2]
		name := bop.byte255ToId(name_b)
		oneIndex := ContextIndex{
			Status:      status,
			ContextName: name,
		}
		index = append(index, oneIndex)
		i = i + 1 + DOT_ID_MAX_LENGTH_V2
	}

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
