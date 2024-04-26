// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ]
// Stephen Fire Meditation Qin [ 火志溟 ] -> stephenfmqin@gmail.com
// This source code is governed by GNU LGPL v3 license

package dot

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	running        bool                // 是否在运行状态

}

// dot的操作锁
type DotLock struct {
	LockTime time.Time
	LockType _BlockDotLockType // 锁的类型
	Lock     *sync.RWMutex
}

// 新建一个dot
func (bop *BlockOp) AddDot(id string, data []byte) (fpath string, fname string, err error) {
	return bop.NewDot(id, data)
}
func (bop *BlockOp) NewDot(id string, data []byte) (fpath string, fname string, err error) {
	if bop.running == false {
		err = fmt.Errorf("The Dot Block is Stop!")
		return
	}

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
	if bop.running == false {
		err = fmt.Errorf("The Dot Block is Stop!")
		return
	}

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
	if bop.running == false {
		err = fmt.Errorf("The Dot Block is Stop!")
		return
	}

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
	if bop.running == false {
		err = fmt.Errorf("The Dot Block is Stop!")
		return
	}

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

	/*
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
				the_c_up_data_name := fpath + fname + "_context_" + base.GetSha1Sum(context_index[cindex].ContextName) + "_UP"
				os.Remove(the_c_up_data_name)
				// TODO Down关系的外部数据文件删除
			}
		}
		// 开始删文件了
		os.Remove(fpath + fname + "_context_index")
		os.Remove(fpath + fname + "_context_del_index")
		os.Remove(fpath + fname + "_body")
	*/

	// 不读，直接删
	files, err := filepath.Glob(fpath + fname + "_*")
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	for _, f := range files {
		if err = os.Remove(f); err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
	}

	return
}

// 增加一个context
//
// 目前暂时不会去补index中早先被删除掉的位置
func (bop *BlockOp) AddContext(dotid string, contextname string) (err error) {
	if bop.running == false {
		err = fmt.Errorf("The Dot Block is Stop!")
		return
	}

	fname, fpath, err := bop.findFilePath(dotid)
	if err != nil {
		err = fmt.Errorf("%v", err)
		return
	}
	if contextname == "" || len([]byte(contextname)) > DOT_ID_MAX_LENGTH_V2 {
		err = fmt.Errorf("Dot Block: The Context name length must less than %v: \"%v\"", DOT_ID_MAX_LENGTH_V2, contextname)
		return
	}
	contextid := base.GetSha1Sum(contextname)

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

	// 看dot存不存在
	ishave := base.FileExist(fpath + fname + "_body")
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the dot \"%v\".", dotid)
		return
	}

	// 看context是否存在，这个简单处理，只看文件是否存在
	ishave = base.FileExist(fpath + fname + "_context_" + contextid)
	// 存在就不再干什么了
	if ishave == true {
		err = fmt.Errorf("Dot Block: The Context \"%v\" is already exist.", contextname)
		return
	}

	// 开始准备空的context
	// 准备基本头部数据
	dotversion_b := iendecode.Uint8ToBytes(DOT_NOW_DEFAULT_VERSION) // dot程序版本
	opversion_b := iendecode.Uint64ToBytes(0)                       // 操作版本
	//打开写入context文件
	context_f, err := os.OpenFile(fpath+fname+"_context_"+contextid, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	defer context_f.Close()
	// 开始写
	_, err = context_f.Write(dotversion_b) // 应用版本
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	_, err = context_f.Write(opversion_b) // 操作版本
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	_, err = context_f.Write(bop.idToByte255(contextname)) // ContextName
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	_, err = context_f.Write(bop.idToByte255("")) // Context的空up关系
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	_, err = context_f.Write(iendecode.Uint8ToBytes(uint8(DOT_CONTEXT_UP_DOWN_INDEX_NOTHING))) // Context的默认配置状态
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	_, err = context_f.Write(iendecode.Uint64ToBytes(uint64(0))) // Context的空up配置数据长度
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	uldata := make([]byte, DOT_CONTENT_MAX_IN_DATA_V2)
	_, err = context_f.Write(uldata) // Context的空up配置数据
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	//打开写入context的down删除索引文件
	dot_context_down_del_index_f, err := os.OpenFile(fpath+fname+"_context_"+contextid+"_del_index", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	defer dot_context_down_del_index_f.Close()
	// 开始写
	_, err = dot_context_down_del_index_f.Write(dotversion_b) // 应用版本
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	_, err = dot_context_down_del_index_f.Write(opversion_b) // 操作版本
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	// 更新索引
	context_index_f, err := os.OpenFile(fpath+fname+"_context_index", os.O_RDWR, 0600)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	defer context_index_f.Close()
	// 获取操作版本，并且+1
	opversion_b = make([]byte, 8)
	read_n, err := context_index_f.ReadAt(opversion_b, 1)
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
	// 写新的操作版本
	_, err = context_index_f.WriteAt(opversion_b, 1)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	// 查找文件末尾的偏移量
	theend, _ := context_index_f.Seek(0, os.SEEK_END)
	// 写入默认的状态
	_, err = context_index_f.WriteAt(iendecode.Uint8ToBytes(uint8(DOT_CONTEXT_INDEX_NOTHING)), theend)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	// 写入上下文关系的名字
	_, err = context_index_f.WriteAt(bop.idToByte255(contextname), theend+1)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	return

}

// 读取所有context名称
func (bop *BlockOp) ReadAllContextName(dotid string) (index []string, err error) {
	return bop.ShowAllContextName(dotid)
}

// 读取所有context名称
func (bop *BlockOp) ShowAllContextName(dotid string) (index []string, err error) {
	if bop.running == false {
		err = fmt.Errorf("The Dot Block is Stop!")
		return
	}

	fname, fpath, err := bop.findFilePath(dotid)
	if err != nil {
		err = fmt.Errorf("%v", err)
		return
	}

	//加读锁
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

	// 看dot存不存在
	ishave := base.FileExist(fpath + fname + "_body")
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the dot \"%v\".", dotid)
		return
	}

	// 读context索引
	context_b, context_b_l, err := bop.readAfter(1+8, fpath+fname+"_context_index")
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	index = make([]string, 0)
	var i int64 = 0
	b_len := int64(context_b_l)
	for {
		if i >= b_len {
			break
		}
		// 状态
		status_b := context_b[i : i+1]
		status_uint := iendecode.BytesToUint8(status_b)
		status := _DotContextIndex_Status(status_uint)
		name_b := context_b[i+1 : i+1+DOT_ID_MAX_LENGTH_V2]
		name := bop.byte255ToId(name_b)
		if status != DOT_CONTEXT_INDEX_DEL {
			index = append(index, name)
		}
		i = i + 1 + DOT_ID_MAX_LENGTH_V2
	}

	return
}

// 删除一个context
func (bop *BlockOp) DelContext(dotid, contextname string) (err error) {
	return bop.DropContext(dotid, contextname)
}
func (bop *BlockOp) DropContext(dotid, contextname string) (err error) {
	if bop.running == false {
		err = fmt.Errorf("The Dot Block is Stop!")
		return
	}

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

	// 看dot存不存在
	ishave := base.FileExist(fpath + fname + "_body")
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the dot \"%v\".", dotid)
		return
	}

	// 看context存在不存在
	contextid := base.GetSha1Sum(contextname)
	ishave = base.FileExist(fpath + fname + "_context_" + contextid)
	if ishave != true {
		// 不存在就不管了
		return
	}

	// 打开context索引
	context_index_f, err := os.OpenFile(fpath+fname+"_context_index", os.O_RDWR, 0600)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	defer context_index_f.Close()
	// 获取操作版本，并且+1
	opversion_b := make([]byte, 8)
	read_n, err := context_index_f.ReadAt(opversion_b, 1)
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

	/*
		// 打开删除索引文件
		context_index_del_f, err := os.OpenFile(fpath+fname+"_context_del_index", os.O_RDWR, 0600)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		defer context_index_del_f.Close()
		// 获取操作版本，并且+1
		del_opversion_b := make([]byte, 8)
		del_read_n, err := context_index_del_f.ReadAt(del_opversion_b, 1)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		if del_read_n != 8 {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		del_opversion := iendecode.BytesToUint64(del_opversion_b)
		del_opversion++
		del_opversion_b = iendecode.Uint64ToBytes(del_opversion)
	*/

	// 遍历索引
	context_b, context_b_l, err := bop.readAfterWithFile(1+8, context_index_f)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	var i int64 = 0        // 字节计数
	var index_i uint64 = 0 // 索引位置计数
	b_len := int64(context_b_l)
	for {
		if i >= b_len {
			break
		}
		// 状态
		status_b := context_b[i : i+1]
		status_uint := iendecode.BytesToUint8(status_b)
		status := _DotContextIndex_Status(status_uint)
		name_b := context_b[i+1 : i+1+DOT_ID_MAX_LENGTH_V2]
		name := bop.byte255ToId(name_b)
		// 如果找到了这个context，并且没有被标记过删除
		if name == contextname && status != DOT_CONTEXT_INDEX_DEL {
			// 去改写状态，标记为删除
			context_index_f.WriteAt(iendecode.Uint8ToBytes(uint8(DOT_CONTEXT_INDEX_DEL)), 1+8+i)
			// 写新的索引操作版本
			_, err = context_index_f.WriteAt(opversion_b, 1)
			if err != nil {
				err = fmt.Errorf("Dot Block: %v", err)
				return
			}

			// 重整删除索引
			err = bop.reAddDelIndexDirect(fpath+fname+"_context_del_index", index_i)
			if err != nil {
				err = fmt.Errorf("Dot Block: %v", err)
				return
			}
			/*
				// 写新的删除索引操作版本
				_, err = context_index_del_f.WriteAt(del_opversion_b, 1)
				if err != nil {
					err = fmt.Errorf("Dot Block: %v", err)
					return
				}
				// 查找删除索引文件末尾的偏移量
				theend, _ := context_index_del_f.Seek(0, os.SEEK_END)
				// 写入index_i索引位置
				_, err = context_index_del_f.WriteAt(iendecode.Uint64ToBytes(index_i), theend)
				if err != nil {
					err = fmt.Errorf("Dot Block: %v", err)
					return
				}
			*/
			// 删除context的相关文件
			var context_files []string
			context_files, err = filepath.Glob(fpath + fname + "_context_" + contextid + "_*")
			if err != nil {
				err = fmt.Errorf("Dot Block: %v", err)
				return
			}
			for _, f := range context_files {
				if err = os.Remove(f); err != nil {
					err = fmt.Errorf("Dot Block: %v", err)
					return
				}
			}
			if err = os.Remove(fpath + fname + "_context_" + contextid); err != nil {
				err = fmt.Errorf("Dot Block: %v", err)
				return
			}

			break // 跳出循环
		}
		i = i + 1 + DOT_ID_MAX_LENGTH_V2
		index_i++
	}

	return
}

// 修改一个context的up信息（名称+数据）
func (bop *BlockOp) UpdateContextUp(dotid, contextname, upname string, updata []byte) (err error) {
	return
}

// 修改一个context的up信息（只名称）
func (bop *BlockOp) UpdateContextUpName(dotid, contextname, upname string) (err error) {
	if bop.running == false {
		err = fmt.Errorf("The Dot Block is Stop!")
		return
	}

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

	// 看dot存不存在
	ishave := base.FileExist(fpath + fname + "_body")
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the dot \"%v\".", dotid)
		return
	}

	// 看context存在不存在
	contextid := base.GetSha1Sum(contextname)
	ishave = base.FileExist(fpath + fname + "_context_" + contextid)
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the context \"%v\" in dot \"%v\"", contextname, dotid)
		return
	}

	// 打开数据文件写入
	context_f, err := os.OpenFile(fpath+fname+"_context_"+contextid, os.O_RDWR, 0600)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	defer context_f.Close()

	// 获取操作版本，并且+1
	opversion_b := make([]byte, 8)
	read_n, err := context_f.ReadAt(opversion_b, 1)
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

	// 编码名字
	upname_b := bop.idToByte255(upname)

	// 写新的索引操作版本
	_, err = context_f.WriteAt(opversion_b, 1)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	// 写入编码名字
	_, err = context_f.WriteAt(upname_b, 1+8+DOT_ID_MAX_LENGTH_V2)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	return
}

// 修改一个context的up信息（只数据）
func (bop *BlockOp) UpdateContextUpData(dotid, contextname string, updata []byte) (err error) {
	if bop.running == false {
		err = fmt.Errorf("The Dot Block is Stop!")
		return
	}

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

	// 看dot存不存在
	ishave := base.FileExist(fpath + fname + "_body")
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the dot \"%v\".", dotid)
		return
	}

	// 看context存在不存在
	contextid := base.GetSha1Sum(contextname)
	ishave = base.FileExist(fpath + fname + "_context_" + contextid)
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the context \"%v\" in dot \"%v\"", contextname, dotid)
		return
	}

	// 打开数据文件写入
	context_f, err := os.OpenFile(fpath+fname+"_context_"+contextid, os.O_RDWR, 0600)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	defer context_f.Close()

	// 获取操作版本，并且+1
	opversion_b := make([]byte, 8)
	read_n, err := context_f.ReadAt(opversion_b, 1)
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

	// 如果存在外部数据文件，先删了再说
	ishave = base.FileExist(fpath + fname + "_context_" + contextid + "_UP")
	if ishave == true {
		if err = os.Remove(fpath + fname + "_context_" + contextid + "_UP"); err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
	}
	// 查看数据长度
	data_len := len(updata)
	if data_len <= DOT_CONTENT_MAX_IN_DATA_V2 {
		// 如果data长度可以
		// 写入数据状态
		_, err = context_f.WriteAt(iendecode.Uint8ToBytes(uint8(DOT_CONTEXT_UP_DOWN_INDEX_INDATA)), 1+8+DOT_ID_MAX_LENGTH_V2+DOT_ID_MAX_LENGTH_V2)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		// 写入数据长度
		_, err = context_f.WriteAt(iendecode.Uint64ToBytes(uint64(data_len)), 1+8+DOT_ID_MAX_LENGTH_V2+DOT_ID_MAX_LENGTH_V2+1)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		// 写入数据
		_, err = context_f.WriteAt(updata, 1+8+DOT_ID_MAX_LENGTH_V2+DOT_ID_MAX_LENGTH_V2+1+8)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
	} else {
		// 如果data长度太长

		// 把内容写入文件
		err = ioutil.WriteFile(fpath+fname+"_context_"+contextid+"_UP", updata, 0600)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		// 写入数据状态
		_, err = context_f.WriteAt(iendecode.Uint8ToBytes(uint8(DOT_CONTEXT_UP_DOWN_INDEX_OUTDATA)), 1+8+DOT_ID_MAX_LENGTH_V2+DOT_ID_MAX_LENGTH_V2)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		// 写入数据长度
		_, err = context_f.WriteAt(iendecode.Uint64ToBytes(uint64(data_len)), 1+8+DOT_ID_MAX_LENGTH_V2+DOT_ID_MAX_LENGTH_V2+1)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
	}

	// 写新的索引操作版本
	_, err = context_f.WriteAt(opversion_b, 1)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	return
}

// 读取一个context的up信息(名称)
func (bop *BlockOp) ReadContextUpName(dotid, contextname string) (upname string, err error) {
	if bop.running == false {
		err = fmt.Errorf("The Dot Block is Stop!")
		return
	}

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

	// 看dot存在否
	fname_data := fname + "_body"
	ishave := base.FileExist(fpath + fname_data)
	// 如果不存在就返回错误
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the dot \"%v\".", dotid)
		return
	}
	// 看context存在不存在
	contextid := base.GetSha1Sum(contextname)
	ishave = base.FileExist(fpath + fname + "_context_" + contextid)
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the context \"%v\" in dot \"%v\"", contextname, dotid)
		return
	}

	// 打开文件
	context_f, err := os.OpenFile(fpath+fname+"_context_"+contextid, os.O_RDONLY, 0600)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	// 获取名字
	upname_b := make([]byte, DOT_ID_MAX_LENGTH_V2)
	read_n, err := context_f.ReadAt(upname_b, 1+8+DOT_ID_MAX_LENGTH_V2)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	if read_n != DOT_ID_MAX_LENGTH_V2 {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	upname = bop.byte255ToId(upname_b)

	return
}

// 读取一个context的up信息(数据)
func (bop *BlockOp) ReadContextUpData(dotid, contextname string) (updata []byte, err error) {
	if bop.running == false {
		err = fmt.Errorf("The Dot Block is Stop!")
		return
	}

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

	// 看dot存在否
	fname_data := fname + "_body"
	ishave := base.FileExist(fpath + fname_data)
	// 如果不存在就返回错误
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the dot \"%v\".", dotid)
		return
	}
	// 看context存在不存在
	contextid := base.GetSha1Sum(contextname)
	ishave = base.FileExist(fpath + fname + "_context_" + contextid)
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the context \"%v\" in dot \"%v\"", contextname, dotid)
		return
	}

	// 打开文件
	context_f, err := os.OpenFile(fpath+fname+"_context_"+contextid, os.O_RDONLY, 0600)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	// 获取data的保存状态
	data_status_b := make([]byte, 1)
	read_n, err := context_f.ReadAt(data_status_b, 1+8+DOT_ID_MAX_LENGTH_V2+DOT_ID_MAX_LENGTH_V2)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	if read_n != 1 {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	data_status := _DotContextUpDownIndex_Status(iendecode.BytesToUint8(data_status_b))
	if data_status == DOT_CONTEXT_UP_DOWN_INDEX_NOTHING {
		return // 没有数据就返回空
	} else if data_status == DOT_CONTEXT_UP_DOWN_INDEX_INDATA {
		// 数据在内部
		// 读数据的长度
		thelen_b := make([]byte, 8)
		read_n, err = context_f.ReadAt(thelen_b, 1+8+DOT_ID_MAX_LENGTH_V2+DOT_ID_MAX_LENGTH_V2+1)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		if read_n != 8 {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		thelen := iendecode.BytesToInt(thelen_b)
		// 读数据
		updata = make([]byte, thelen)
		read_n, err = context_f.ReadAt(updata, 1+8+DOT_ID_MAX_LENGTH_V2+DOT_ID_MAX_LENGTH_V2+1+8)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		if read_n != thelen {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
	} else if data_status == DOT_CONTEXT_UP_DOWN_INDEX_OUTDATA {
		// 数据在外部，就直接读整个文件
		updata, err = ioutil.ReadFile(fpath + fname + "_context_" + contextid + "_UP")
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
	} else {
		return
	}

	return
}

// 增加一个context的down信息（名称+数据）
//
// 如果有同名但被删除的位置，则去替代这个位置
//
// 如果有删除的空位（_del_index中记录），则占用记录中的第一个位置
func (bop *BlockOp) AddContextDown(dotid, contextname, downname string, data []byte) (err error) {
	if bop.running == false {
		err = fmt.Errorf("The Dot Block is Stop!")
		return
	}

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

	// 看dot存不存在
	ishave := base.FileExist(fpath + fname + "_body")
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the dot \"%v\".", dotid)
		return
	}

	// 看context存在不存在
	contextid := base.GetSha1Sum(contextname)
	ishave = base.FileExist(fpath + fname + "_context_" + contextid)
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the context \"%v\" in dot \"%v\"", contextname, dotid)
		return
	}
	// 看downname是否符合
	if downname == "" || len([]byte(downname)) > DOT_ID_MAX_LENGTH_V2 {
		err = fmt.Errorf("Dot Block: The Context Down name length must less than %v: \"%v\"", DOT_ID_MAX_LENGTH_V2, downname)
		return
	}
	downid := base.GetSha1Sum(downname)

	// 获取down的删除索引
	del_list, del_version, err := bop.readDelList(fpath + fname + "_context_" + contextid + "_del_index")
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	// 打开Context文件的写入
	context_f, err := os.OpenFile(fpath+fname+"_context_"+contextid, os.O_RDWR, 0600)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	defer context_f.Close()

	// 获取所有down的信息
	down_status_list, err := bop.readContextDownStatus(context_f)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	// 查看down是否存在
	exist, down_status := bop.ifDownExist(downname, down_status_list)
	if exist == true && down_status.Status != DOT_CONTEXT_UP_DOWN_INDEX_DEL {
		err = fmt.Errorf("Dot Block: The Context Down is exist: %v", downname)
		return
	}

	// 获取操作版本，并且+1
	opversion_b := make([]byte, 8)
	read_n, err := context_f.ReadAt(opversion_b, 1)
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

	var write_i int64 // 写入位置
	if exist == false {
		// 如果彻底没有这个down
		if len(del_list) == 0 {
			//如果也没有被删除的空位置，就把写入位置调到文件末尾
			write_i, _ = context_f.Seek(0, os.SEEK_END)
		} else {
			//否则选第一个空位置
			write_i = 1 + 8 + DOT_ID_MAX_LENGTH_V2 + DOT_ID_MAX_LENGTH_V2 + 1 + 8 + DOT_CONTENT_MAX_IN_DATA_V2
			write_i = write_i + (int64(del_list[0]) * (1 + DOT_ID_MAX_LENGTH_V2 + 8 + DOT_CONTENT_MAX_IN_DATA_V2))
			// 并且去更新删除索引
			err = bop.reDelIndex(fpath+fname+"_context_"+contextid+"_del_index", del_list, del_list[0], del_version+1)
			if err != nil {
				err = fmt.Errorf("Dot Block: %v", err)
				return
			}
		}

	} else {
		// 如果有同名，但被删除了，则把写入位置调整到这个位置
		write_i = int64(down_status.HardCount)
		// 并且去更新删除索引
		err = bop.reDelIndex(fpath+fname+"_context_"+contextid+"_del_index", del_list, down_status.HardIndex, del_version+1)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
	}

	// 如果存在外部数据文件，先删了再说

	ishave = base.FileExist(fpath + fname + "_context_" + contextid + "_DOWN_" + downid)
	if ishave == true {
		if err = os.Remove(fpath + fname + "_context_" + contextid + "_DOWN_" + downid); err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
	}

	// 查看数据长度
	data_len := len(data)
	if data_len <= DOT_CONTENT_MAX_IN_DATA_V2 {
		// 如果data长度可以
		// 写入数据状态
		_, err = context_f.WriteAt(iendecode.Uint8ToBytes(uint8(DOT_CONTEXT_UP_DOWN_INDEX_INDATA)), write_i)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		// 写入down的名字
		downname_b := bop.idToByte255(downname)
		_, err = context_f.WriteAt(downname_b, write_i+1)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		// 写入数据长度
		_, err = context_f.WriteAt(iendecode.Uint64ToBytes(uint64(data_len)), write_i+1+DOT_ID_MAX_LENGTH_V2)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		// 先写空的数据占位置
		uldata := make([]byte, DOT_CONTENT_MAX_IN_DATA_V2)
		_, err = context_f.WriteAt(uldata, write_i+1+DOT_ID_MAX_LENGTH_V2+8)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		// 写入真实数据
		_, err = context_f.WriteAt(data, write_i+1+DOT_ID_MAX_LENGTH_V2+8)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
	} else {
		// 如果data长度太长

		// 把内容写入文件
		err = ioutil.WriteFile(fpath+fname+"_context_"+contextid+"_DOWN_"+downid, data, 0600)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}

		// 写入数据状态
		_, err = context_f.WriteAt(iendecode.Uint8ToBytes(uint8(DOT_CONTEXT_UP_DOWN_INDEX_OUTDATA)), write_i)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		// 写入down的名字
		downname_b := bop.idToByte255(downname)
		_, err = context_f.WriteAt(downname_b, write_i+1)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		// 写入数据长度
		_, err = context_f.WriteAt(iendecode.Uint64ToBytes(uint64(data_len)), write_i+1+DOT_ID_MAX_LENGTH_V2)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		// 写空的数据占位置
		uldata := make([]byte, DOT_CONTENT_MAX_IN_DATA_V2)
		_, err = context_f.WriteAt(uldata, write_i+1+DOT_ID_MAX_LENGTH_V2+8)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
	}

	// 写新的操作版本
	_, err = context_f.WriteAt(opversion_b, 1)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	return
}

// 读所有的Down关系名称
func (bop *BlockOp) ShowContextAllDownName(dotid, contextname string) (index []string, err error) {
	return bop.ReadContextAllDownName(dotid, contextname)
}

// 读所有的Down关系名称
func (bop *BlockOp) ReadContextAllDownName(dotid, contextname string) (index []string, err error) {
	if bop.running == false {
		err = fmt.Errorf("The Dot Block is Stop!")
		return
	}

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

	// 看dot存在否
	fname_data := fname + "_body"
	ishave := base.FileExist(fpath + fname_data)
	// 如果不存在就返回错误
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the dot \"%v\".", dotid)
		return
	}
	// 看context存在不存在
	contextid := base.GetSha1Sum(contextname)
	ishave = base.FileExist(fpath + fname + "_context_" + contextid)
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the context \"%v\" in dot \"%v\"", contextname, dotid)
		return
	}

	// 打开文件
	context_f, err := os.OpenFile(fpath+fname+"_context_"+contextid, os.O_RDONLY, 0600)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	// 读所有Down的状态
	down_status, err := bop.readContextDownStatus(context_f)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	index = make([]string, 0)
	for _, one := range down_status {
		if one.Status != DOT_CONTEXT_UP_DOWN_INDEX_DEL {
			index = append(index, one.Name)
		}
	}

	return
}

// 读一个down的数据
func (bop *BlockOp) ReadContextOneDownData(dotid, contextname, downname string) (data []byte, err error) {
	if bop.running == false {
		err = fmt.Errorf("The Dot Block is Stop!")
		return
	}

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

	// 看dot存在否
	fname_data := fname + "_body"
	ishave := base.FileExist(fpath + fname_data)
	// 如果不存在就返回错误
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the dot \"%v\".", dotid)
		return
	}
	// 看context存在不存在
	contextid := base.GetSha1Sum(contextname)
	ishave = base.FileExist(fpath + fname + "_context_" + contextid)
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the context \"%v\" in dot \"%v\"", contextname, dotid)
		return
	}

	// 打开文件
	context_f, err := os.OpenFile(fpath+fname+"_context_"+contextid, os.O_RDONLY, 0600)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	defer context_f.Close()

	// 读所有Down的状态
	down_status, err := bop.readContextDownStatus(context_f)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	// 看down存在不存在
	exist, status := bop.ifDownExist(downname, down_status)
	if exist == false || status.Status == DOT_CONTEXT_UP_DOWN_INDEX_DEL {
		err = fmt.Errorf("Dot Block: Can not find the Context Down \"%v\".", downname)
		return
	}
	// 看是内部外部数据
	if status.Status == DOT_CONTEXT_UP_DOWN_INDEX_INDATA {
		// 读数据的长度
		thelen_b := make([]byte, 8)
		read_n := int(0)
		read_n, err = context_f.ReadAt(thelen_b, int64(status.HardCount+1+DOT_ID_MAX_LENGTH_V2))
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		if read_n != 8 {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		thelen := iendecode.BytesToInt(thelen_b)
		// 读数据
		data = make([]byte, thelen)
		read_n, err = context_f.ReadAt(data, int64(status.HardCount+1+DOT_ID_MAX_LENGTH_V2+8))
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		if read_n != thelen {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
	} else if status.Status == DOT_CONTEXT_UP_DOWN_INDEX_OUTDATA {
		data, err = ioutil.ReadFile(fpath + fname + "_context_" + contextid + "_DOWN_" + base.GetSha1Sum(downname))
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
	} else {
		err = fmt.Errorf("Dot Block: Can not find the Context Down \"%v\".", downname)
		return
	}

	return
}

// 修改一个context的down信息（只数据）
func (bop *BlockOp) UpdateContextDownData(dotid, contextname, downname string, data []byte) (err error) {
	if bop.running == false {
		err = fmt.Errorf("The Dot Block is Stop!")
		return
	}

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

	// 看dot存不存在
	ishave := base.FileExist(fpath + fname + "_body")
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the dot \"%v\".", dotid)
		return
	}

	// 看context存在不存在
	contextid := base.GetSha1Sum(contextname)
	ishave = base.FileExist(fpath + fname + "_context_" + contextid)
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the context \"%v\" in dot \"%v\"", contextname, dotid)
		return
	}

	// 打开context文件写入
	context_f, err := os.OpenFile(fpath+fname+"_context_"+contextid, os.O_RDWR, 0600)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	defer context_f.Close()

	// 获取down表
	down_list, err := bop.readContextDownStatus(context_f)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	// 查看down是否存在
	exist, downstatus := bop.ifDownExist(downname, down_list)
	if exist == false || downstatus.Status == DOT_CONTEXT_UP_DOWN_INDEX_DEL {
		err = fmt.Errorf("Dot Block: The Context Down is not exist: \"%v\"", downname)
		return
	}

	// 获取downid
	downid := base.GetSha1Sum(downname)

	// 获取操作版本，并且+1
	opversion_b := make([]byte, 8)
	read_n, err := context_f.ReadAt(opversion_b, 1)
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

	// 如果存在外部数据文件，先删了再说
	ishave = base.FileExist(fpath + fname + "_context_" + contextid + "_DOWN_" + downid)
	if ishave == true {
		if err = os.Remove(fpath + fname + "_context_" + contextid + "_DOWN_" + downid); err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
	}
	// 查看数据长度
	data_len := len(data)
	if data_len <= DOT_CONTENT_MAX_IN_DATA_V2 {
		// 如果data长度可以
		// 写入数据状态
		_, err = context_f.WriteAt(iendecode.Uint8ToBytes(uint8(DOT_CONTEXT_UP_DOWN_INDEX_INDATA)), int64(downstatus.HardCount))
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		// 写入数据长度
		_, err = context_f.WriteAt(iendecode.Uint64ToBytes(uint64(data_len)), int64(downstatus.HardCount+1+DOT_ID_MAX_LENGTH_V2))
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		// 清空位置
		uldata := make([]byte, DOT_CONTENT_MAX_IN_DATA_V2)
		_, err = context_f.WriteAt(uldata, int64(downstatus.HardCount+1+DOT_ID_MAX_LENGTH_V2+8))
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		// 写入数据
		_, err = context_f.WriteAt(data, int64(downstatus.HardCount+1+DOT_ID_MAX_LENGTH_V2+8))
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
	} else {
		// 如果data长度太长

		// 把内容写入文件
		err = ioutil.WriteFile(fpath+fname+"_context_"+contextid+"_DOWN_"+downid, data, 0600)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		// 写入数据状态
		_, err = context_f.WriteAt(iendecode.Uint8ToBytes(uint8(DOT_CONTEXT_UP_DOWN_INDEX_OUTDATA)), int64(downstatus.HardCount))
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		// 写入数据长度
		_, err = context_f.WriteAt(iendecode.Uint64ToBytes(uint64(data_len)), int64(downstatus.HardCount+1+DOT_ID_MAX_LENGTH_V2))
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
	}

	// 写新的索引操作版本
	_, err = context_f.WriteAt(opversion_b, 1)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	return
}

// 删除一个context的down信息
func (bop *BlockOp) DropContextDownData(dotid, contextname, downname string) (err error) {
	return bop.DelContextDownData(dotid, contextname, downname)
}

// 删除一个context的down信息
func (bop *BlockOp) DelContextDownData(dotid, contextname, downname string) (err error) {
	if bop.running == false {
		err = fmt.Errorf("The Dot Block is Stop!")
		return
	}

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

	// 看dot存不存在
	ishave := base.FileExist(fpath + fname + "_body")
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the dot \"%v\".", dotid)
		return
	}

	// 看context存在不存在
	contextid := base.GetSha1Sum(contextname)
	ishave = base.FileExist(fpath + fname + "_context_" + contextid)
	if ishave != true {
		err = fmt.Errorf("Dot Block: Can not find the context \"%v\" in dot \"%v\"", contextname, dotid)
		return
	}

	// 打开context文件写入
	context_f, err := os.OpenFile(fpath+fname+"_context_"+contextid, os.O_RDWR, 0600)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	defer context_f.Close()

	// 获取down表
	down_list, err := bop.readContextDownStatus(context_f)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	// 查看down是否存在
	exist, downstatus := bop.ifDownExist(downname, down_list)
	if exist == false || downstatus.Status == DOT_CONTEXT_UP_DOWN_INDEX_DEL {
		// 没有就返回，什么事情都不做
		return
	}

	// 获取downid
	downid := base.GetSha1Sum(downname)

	// 获取操作版本，并且+1
	opversion_b := make([]byte, 8)
	read_n, err := context_f.ReadAt(opversion_b, 1)
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

	// 如果存在外部数据文件，先删了再说
	ishave = base.FileExist(fpath + fname + "_context_" + contextid + "_DOWN_" + downid)
	if ishave == true {
		if err = os.Remove(fpath + fname + "_context_" + contextid + "_DOWN_" + downid); err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
	}

	// 添加删除索引
	err = bop.reAddDelIndexDirect(fpath+fname+"_context_"+contextid+"_del_index", downstatus.HardIndex)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	// 修改down的删除状态
	_, err = context_f.WriteAt(iendecode.Uint8ToBytes(uint8(DOT_CONTEXT_UP_DOWN_INDEX_DEL)), int64(downstatus.HardCount))
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	// 写新的索引操作版本
	_, err = context_f.WriteAt(opversion_b, 1)
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
func (bop *BlockOp) OutLock(dotid string) (err error) {
	if bop.running == false {
		err = fmt.Errorf("The Dot Block is Stop!")
		return
	}

	_, _, err = bop.findFilePath(dotid)
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
	// 如果没有锁就加外部锁，如果是外部锁，就不管了
	if bop.dots_lock[dotid].LockType == BLOCK_DOT_LOCK_TYPE_NOTHING {
		bop.dots_lock[dotid].LockTime = time.Now()
		bop.dots_lock[dotid].LockType = BLOCK_DOT_LOCK_TYPE_OUTSIDE
		bop.dots_lock[dotid].Lock.Lock()
	}
	bop.dots_lock_lock.Unlock()

	return
}

// 外部加读锁
func (bop *BlockOp) OutRLock(dotid string) (err error) {
	if bop.running == false {
		err = fmt.Errorf("The Dot Block is Stop!")
		return
	}

	_, _, err = bop.findFilePath(dotid)
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
	// 如果没有锁就加外部锁，如果是外部锁，就不管了
	if bop.dots_lock[dotid].LockType == BLOCK_DOT_LOCK_TYPE_NOTHING {
		bop.dots_lock[dotid].LockTime = time.Now()
		bop.dots_lock[dotid].LockType = BLOCK_DOT_LOCK_TYPE_OUTSIDE
		bop.dots_lock[dotid].Lock.RLock()
	}
	bop.dots_lock_lock.Unlock()

	return

}

// 外部解锁
func (bop *BlockOp) OutUnlock(dotid string) (err error) {
	if bop.running == false {
		err = fmt.Errorf("The Dot Block is Stop!")
		return
	}

	_, _, err = bop.findFilePath(dotid)
	if err != nil {
		err = fmt.Errorf("%v", err)
		return
	}

	//加锁
	bop.dots_lock_lock.Lock()
	defer bop.dots_lock_lock.Unlock()

	if _, have := bop.dots_lock[dotid]; have != true {
		return
	}
	// 如果是外部锁
	if bop.dots_lock[dotid].LockType == BLOCK_DOT_LOCK_TYPE_OUTSIDE {
		bop.dots_lock[dotid].Lock.Unlock()
		bop.dots_lock[dotid].LockType = BLOCK_DOT_LOCK_TYPE_NOTHING
	} else {
		err = fmt.Errorf("This is not a outside lock!")
		return
	}

	return
}

// 外部读解锁
func (bop *BlockOp) OutRUnlock(dotid string) (err error) {
	if bop.running == false {
		err = fmt.Errorf("The Dot Block is Stop!")
		return
	}

	_, _, err = bop.findFilePath(dotid)
	if err != nil {
		err = fmt.Errorf("%v", err)
		return
	}

	//加锁
	bop.dots_lock_lock.Lock()
	defer bop.dots_lock_lock.Unlock()

	if _, have := bop.dots_lock[dotid]; have != true {
		return
	}
	// 如果是外部锁
	if bop.dots_lock[dotid].LockType == BLOCK_DOT_LOCK_TYPE_OUTSIDE {
		bop.dots_lock[dotid].Lock.RUnlock()
		bop.dots_lock[dotid].LockType = BLOCK_DOT_LOCK_TYPE_NOTHING
	} else {
		err = fmt.Errorf("This is not a outside lock!")
		return
	}

	return

}

// 直接追加删除索引
func (bop *BlockOp) reAddDelIndexDirect(fname string, i uint64) (err error) {
	index, opversion, err := bop.readDelList(fname)
	if err != nil {
		return
	}
	err = bop.reAddDelIndex(fname, index, i, opversion+1)
	return
}

// 重新整理删除索引，加上一个被删除的
func (bop *BlockOp) reAddDelIndex(fname string, index []uint64, i uint64, opversion uint64) (err error) {
	new_index := append(index, i)

	new_index_b, err := iendecode.SliceToBytes("[]uint64", new_index)
	if err != nil {
		return
	}

	// 把老的文件删了
	if err = os.Remove(fname); err != nil {
		return
	}
	//打开写入新的删除索引文件
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	// 开始写
	dotversion_b := iendecode.Uint8ToBytes(DOT_NOW_DEFAULT_VERSION) // dot程序版本
	opversion_b := iendecode.Uint64ToBytes(opversion)               // 操作版本
	_, err = f.Write(dotversion_b)                                  // 应用版本
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	_, err = f.Write(opversion_b) // 操作版本
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	_, err = f.Write(new_index_b) // 索引
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	return
}

// 重新整理删除索引，去除不是删除状态的
func (bop *BlockOp) reDelIndex(fname string, index []uint64, i uint64, opversion uint64) (err error) {
	new_index := make([]uint64, 0)
	for _, one := range index {
		if one != i {
			new_index = append(new_index, i)
		}
	}
	new_index_b, err := iendecode.SliceToBytes("[]uint64", new_index)
	if err != nil {
		return
	}

	// 把老的文件删了
	if err = os.Remove(fname); err != nil {
		return
	}
	//打开写入新的删除索引文件
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	// 开始写
	dotversion_b := iendecode.Uint8ToBytes(DOT_NOW_DEFAULT_VERSION) // dot程序版本
	opversion_b := iendecode.Uint64ToBytes(opversion)               // 操作版本
	_, err = f.Write(dotversion_b)                                  // 应用版本
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	_, err = f.Write(opversion_b) // 操作版本
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	_, err = f.Write(new_index_b) // 索引
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	return
}

// 查看down是否存在
func (bop *BlockOp) ifDownExist(downname string, list []ContextDownStatus) (exist bool, status ContextDownStatus) {
	exist = false
	for _, one := range list {
		if one.Name == downname {
			exist = true
			status = one
			return
		}
	}

	return
}

// 返回完整的被删除列表，直接打开文件全读出来
func (bop *BlockOp) readDelList(fname string) (list []uint64, opversion uint64, err error) {
	f_b, err := ioutil.ReadFile(fname)
	if err != nil {
		return
	}
	opversion_b := f_b[1 : 1+8]
	opversion = iendecode.BytesToUint64(opversion_b)
	list = make([]uint64, 0)
	if len(f_b) == 0 {
		return
	}
	if len(f_b[1+8:]) == 0 {
		return
	}
	list_i, err := iendecode.BytesToSlice("[]uint64", f_b[1+8:])
	if err != nil {
		//err = fmt.Errorf("22 %v", err)
		return
	}
	list = list_i.([]uint64)

	return
}

// 返回完整的Context内Down关系状态索引
func (bop *BlockOp) readContextDownStatus(f *os.File) (d_index []ContextDownStatus, err error) {
	d_index = make([]ContextDownStatus, 0)
	var hard_index uint64 = 0 // 物理索引位置
	f_status, err := f.Stat()
	if err != nil {
		return
	}

	var i int64 = 1 + 8 + DOT_ID_MAX_LENGTH_V2 + DOT_ID_MAX_LENGTH_V2 + 1 + 8 + DOT_CONTENT_MAX_IN_DATA_V2 // 字节计数，从down开始
	b_len := f_status.Size()                                                                               // 数据总长
	down_len := b_len - (1 + 8 + DOT_ID_MAX_LENGTH_V2 + DOT_ID_MAX_LENGTH_V2 + 1 + 8 + DOT_CONTENT_MAX_IN_DATA_V2)

	// 看长度是否正确
	if down_len != 0 && down_len%(1+DOT_ID_MAX_LENGTH_V2+8+DOT_CONTENT_MAX_IN_DATA_V2) != 0 {
		err = fmt.Errorf("The Data's length is wrong.")
		return
	}

	for {
		if i >= b_len {
			break
		}
		var read_n int
		// 获取状态
		status_b := make([]byte, 1)
		read_n, err = f.ReadAt(status_b, i)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		if read_n != 1 {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		status_uint := iendecode.BytesToUint8(status_b)
		status := _DotContextUpDownIndex_Status(status_uint)
		// 获取名字
		name_b := make([]byte, DOT_ID_MAX_LENGTH_V2)
		read_n, err = f.ReadAt(name_b, i+1)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		if read_n != DOT_ID_MAX_LENGTH_V2 {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		name := bop.byte255ToId(name_b)
		// 获取data长度
		dlen_b := make([]byte, 8)
		read_n, err = f.ReadAt(dlen_b, i+1+DOT_ID_MAX_LENGTH_V2)
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		if read_n != 8 {
			err = fmt.Errorf("Dot Block: %v", err)
			return
		}
		dlen := iendecode.BytesToUint64(dlen_b)

		one := ContextDownStatus{
			HardIndex: hard_index,
			HardCount: uint64(i),
			Name:      name,
			Status:    status,
			DataLen:   dlen,
		}
		d_index = append(d_index, one)

		i = i + 1 + DOT_ID_MAX_LENGTH_V2 + 8 + DOT_CONTENT_MAX_IN_DATA_V2
		hard_index++
	}

	return
}

// 完整读取context的索引
func (bop *BlockOp) readContextIndex(b []byte) (index []ContextIndex) {
	index = make([]ContextIndex, 0)
	var i int64 = 0
	b_len := int64(len(b))
	for {
		if i >= b_len {
			break
		}
		// 状态
		status_b := b[i : i+1]
		status_uint := iendecode.BytesToUint8(status_b)
		status := _DotContextIndex_Status(status_uint)
		name_b := b[i+1 : i+1+DOT_ID_MAX_LENGTH_V2]
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
	if id == "" || len([]byte(id)) > DOT_ID_MAX_LENGTH_V2 {
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
