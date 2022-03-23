// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
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

// 添加一组上下文关系
func (dop *DotsOp) AddContextT(dotid string, contextid string, context *Context, optime time.Time) (err error) {
	fname, fpath, err := dop.findFilePath(dotid)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	fmt.Println(fpath)

	// 加锁
	dop.dots_lock_lock.Lock()
	if _, have := dop.dots_lock[fname]; have != true {
		dop.dots_lock[fname] = new(sync.RWMutex)
	}
	dop.dots_lock[fname].Lock()
	dop.dots_lock_lock.Unlock()
	// 函数退出的解锁
	defer func() {
		dop.dots_lock_lock.Lock()
		dop.dots_lock[fname].Unlock()
		delete(dop.dots_lock, fname)
		dop.dots_lock_lock.Unlock()
	}()

	// 找有没有文件
	fname_one_context := fname + DOT_FILE_NAME_CONTEXT + "_" + base.GetSha1Sum(contextid)
	ishave := base.FileExist(fpath + fname_one_context)
	if ishave == true {
		err = fmt.Errorf("dot: The context \"%v\" in dot \"%v\" is already have.", contextid, dotid)
		return
	}
	fname_context := fname + DOT_FILE_NAME_CONTEXT

	// 打开context索引文件
	f_c_i, err := os.OpenFile(fpath+fname_context, os.O_RDWR, 0600)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	defer f_c_i.Close()

	// 构建当前时间
	optime_b, _ := optime.MarshalBinary()

	var read_n int
	// 获取操作时间
	old_optime_b := make([]byte, 15)
	read_n, err = f_c_i.ReadAt(old_optime_b, 1)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	if read_n != 15 {
		err = fmt.Errorf("dot: not the leght %v")
		return
	}
	old_optime := time.Time{}
	err = old_optime.UnmarshalBinary(old_optime_b)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	// 比较操作时间
	if optime.UnixNano() <= old_optime.UnixNano() {
		err = fmt.Errorf("dot: the opreate lock is wrong for dot %v.", dotid)
		return
	}

	// 获取操作版本
	opversion_b := make([]byte, 8)
	read_n, err = f_c_i.ReadAt(opversion_b, 1+15)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	if read_n != 8 {
		err = fmt.Errorf("dot: net the leght %v")
		return
	}
	opversion := iendecode.BytesToUint64(opversion_b)
	opversion++
	opversion_b = iendecode.Uint64ToBytes(opversion)

	//获取context索引
	context_i_b, context_i_b_l, err := dop.readAfterWithFile(1+15+8, f_c_i)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	var context_s = make([]string, 0)
	if context_i_b_l != 0 {
		var context_i interface{}
		context_i, err = iendecode.BytesToSlice("[]string", context_i_b)
		if err != nil {
			err = fmt.Errorf("dot: %v", err)
			return
		}
		context_s = context_i.([]string)
	}
	context_s = append(context_s, contextid)
	context_s_b, err := iendecode.SliceToBytes("[]string", context_s)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	// 准备写入数据
	w_b := make([]byte, 0)
	w_b = append(w_b, optime_b...)
	w_b = append(w_b, opversion_b...)
	w_b = append(w_b, context_s_b...)

	// 准备单个context文件
	fname_context_this := fname + DOT_FILE_NAME_CONTEXT + "_" + base.GetSha1Sum(contextid)
	dop_context_this_f, err := os.OpenFile(fpath+fname_context_this, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	defer dop_context_this_f.Close()

	// 开始写context_this文件
	context_b, err := context.MarshalBinary()
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	dotversion_b := iendecode.Uint8ToBytes(DOT_NOW_DEFAULT_VERSION) // dot程序版本
	c_opversion_b := iendecode.Uint64ToBytes(1)                     // 操作版本

	_, err = dop_context_this_f.Write(dotversion_b) // 应用版本
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_context_this_f.Write(dop.idToByte255(contextid)) // Context ID
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_context_this_f.Write(optime_b) // 时间
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_context_this_f.Write(c_opversion_b) // 操作版本
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_context_this_f.Write(context_b) // 上下文
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}

	// 写context索引
	err = f_c_i.Truncate(1 + 15 + 8)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = f_c_i.WriteAt(w_b, 1)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}

	return
}

// 删除一组上下文
func (dop *DotsOp) DelContextT(dotid string, contextid string, optime time.Time) (err error) {
	fname, fpath, err := dop.findFilePath(dotid)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	fmt.Println(fpath)

	// 加锁
	dop.dots_lock_lock.Lock()
	if _, have := dop.dots_lock[fname]; have != true {
		dop.dots_lock[fname] = new(sync.RWMutex)
	}
	dop.dots_lock[fname].Lock()
	dop.dots_lock_lock.Unlock()
	// 函数退出的解锁
	defer func() {
		dop.dots_lock_lock.Lock()
		dop.dots_lock[fname].Unlock()
		delete(dop.dots_lock, fname)
		dop.dots_lock_lock.Unlock()
	}()

	// 找有没有文件
	fname_one_context := fname + DOT_FILE_NAME_CONTEXT + "_" + base.GetSha1Sum(contextid)
	fname_context := fname + DOT_FILE_NAME_CONTEXT
	ishave := base.FileExist(fpath + fname_one_context)
	if ishave != true {
		err = fmt.Errorf("dot: Can not find the dot \"%v\", or can not find the context \"%v\".", dotid, contextid)
		return
	}

	// 打开文件
	f_con, err := os.OpenFile(fpath+fname_context, os.O_RDWR, 0600)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	defer f_con.Close()

	// 构建当前操作时间
	optime_b, _ := optime.MarshalBinary() // 操作时间，15的长度

	var read_n int
	// 获取操作时间
	old_optime_b := make([]byte, 15)
	read_n, err = f_con.ReadAt(old_optime_b, 1)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	if read_n != 15 {
		err = fmt.Errorf("dot: not the leght %v")
		return
	}
	old_optime := time.Time{}
	err = old_optime.UnmarshalBinary(old_optime_b)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	// 比较操作时间
	if optime.UnixNano() <= old_optime.UnixNano() {
		err = fmt.Errorf("dot: the opreate lock is wrong for dot %v.", dotid)
		return
	}

	// 获取操作版本
	opversion_b := make([]byte, 8)
	read_n, err = f_con.ReadAt(opversion_b, 1+15)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	if read_n != 8 {
		err = fmt.Errorf("dot: net the leght %v")
		return
	}
	opversion := iendecode.BytesToUint64(opversion_b)
	opversion++
	opversion_b = iendecode.Uint64ToBytes(opversion)

	// 读取context索引
	context_b, context_b_l, err := dop.readAfterWithFile(1+15+8, f_con)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	if context_b_l != 0 {
		var context_i interface{}
		context_i, err = iendecode.BytesToSlice("[]string", context_b)
		if err != nil {
			err = fmt.Errorf("dot: %v", err)
			return
		}
		context_s := context_i.([]string)
		var context_n = make([]string, 0)
		for i := range context_s {
			if context_s[i] != contextid {
				context_n = append(context_n, context_s[i])
			}
		}
		var context_n_b []byte
		context_n_b, err = iendecode.SliceToBytes("[]string", context_n)
		if err != nil {
			err = fmt.Errorf("dot: %v", err)
			return
		}
		// 准备写入数据
		w_b := make([]byte, 0)
		w_b = append(w_b, optime_b...)
		w_b = append(w_b, opversion_b...)
		w_b = append(w_b, context_n_b...)
		// 扔掉之前的数据部分
		err = f_con.Truncate(1 + 15 + 8)
		if err != nil {
			err = fmt.Errorf("dot: %v", err)
			return
		}
		// 重新写
		_, err = f_con.WriteAt(w_b, 1)
		if err != nil {
			err = fmt.Errorf("dot: %v", err)
			return
		}
	}

	err = os.Remove(fpath + fname_one_context)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}

	return
}

// 删除一个上下文中的Down关系
func (dop *DotsOp) DelOneDownT(dotid string, contextid string, downname string, optime time.Time) (err error) {
	fname, fpath, err := dop.findFilePath(dotid)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	fmt.Println(fpath)

	// 加锁
	dop.dots_lock_lock.Lock()
	if _, have := dop.dots_lock[fname]; have != true {
		dop.dots_lock[fname] = new(sync.RWMutex)
	}
	dop.dots_lock[fname].Lock()
	dop.dots_lock_lock.Unlock()
	// 函数退出的解锁
	defer func() {
		dop.dots_lock_lock.Lock()
		dop.dots_lock[fname].Unlock()
		delete(dop.dots_lock, fname)
		dop.dots_lock_lock.Unlock()
	}()

	// 找有没有文件
	fname_one_context := fname + DOT_FILE_NAME_CONTEXT + "_" + base.GetSha1Sum(contextid)
	ishave := base.FileExist(fpath + fname_one_context)
	if ishave != true {
		err = fmt.Errorf("dot: Can not find the dot \"%v\", or can not find the context \"%v\".", dotid, contextid)
		return
	}

	// 打开文件
	f, err := os.OpenFile(fpath+fname_one_context, os.O_RDWR, 0600)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	defer f.Close()

	// 构建当前操作时间
	optime_b, _ := optime.MarshalBinary() // 操作时间，15的长度

	var read_n int
	// 获取操作时间
	old_optime_b := make([]byte, 15)
	read_n, err = f.ReadAt(old_optime_b, 1+DOT_ID_MAX_LENGTH_V1)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	if read_n != 15 {
		err = fmt.Errorf("dot: not the leght %v")
		return
	}
	old_optime := time.Time{}
	err = old_optime.UnmarshalBinary(old_optime_b)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	// 比较操作时间
	if optime.UnixNano() <= old_optime.UnixNano() {
		err = fmt.Errorf("dot: the opreate lock is wrong for dot %v.", dotid)
		return
	}

	// 获取操作版本
	opversion_b := make([]byte, 8)
	read_n, err = f.ReadAt(opversion_b, 1+DOT_ID_MAX_LENGTH_V1+15)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	if read_n != 8 {
		err = fmt.Errorf("dot: net the leght %v")
		return
	}
	opversion := iendecode.BytesToUint64(opversion_b)
	opversion++
	opversion_b = iendecode.Uint64ToBytes(opversion)

	// 读取这个context的结构体
	con_b, len, err := dop.readAfterWithFile(1+DOT_ID_MAX_LENGTH_V1+15+8, f)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	if len == 0 {
		err = fmt.Errorf("dot: The context have some error.")
		return
	}
	context := &Context{}
	err = context.UnmarshalBinary(con_b)

	// 修改
	delete(context.Down, downname)
	context_b, err := context.MarshalBinary()
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}

	// 准备写入数据
	w_b := make([]byte, 0)
	w_b = append(w_b, optime_b...)
	w_b = append(w_b, opversion_b...)
	w_b = append(w_b, context_b...)
	// 扔掉之前的数据部分
	err = f.Truncate(1 + DOT_ID_MAX_LENGTH_V1 + 15 + 8)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	// 重新写
	_, err = f.WriteAt(w_b, 1+DOT_ID_MAX_LENGTH_V1)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}

	return
}

// 新建一个只有数据的dot
func (dop *DotsOp) NewDotT(id string, data []byte, optime time.Time) (err error) {
	fname, fpath, err := dop.findFilePath(id)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	fmt.Println(fpath)

	// 加锁
	dop.dots_lock_lock.Lock()
	if _, have := dop.dots_lock[fname]; have != true {
		dop.dots_lock[fname] = new(sync.RWMutex)
	}
	dop.dots_lock[fname].Lock()
	dop.dots_lock_lock.Unlock()
	// 函数退出的解锁
	defer func() {
		dop.dots_lock_lock.Lock()
		dop.dots_lock[fname].Unlock()
		delete(dop.dots_lock, fname)
		dop.dots_lock_lock.Unlock()
	}()

	fname_data := fname + DOT_FILE_NAME_DATA
	fname_context := fname + DOT_FILE_NAME_CONTEXT

	ishave := base.FileExist(fpath + fname_data)
	if ishave == true {
		err = fmt.Errorf("dot: The dot id \"%v\" already have.", id)
		return
	}

	optime_b, _ := optime.MarshalBinary()                           // 操作时间，15的长度
	dotversion_b := iendecode.Uint8ToBytes(DOT_NOW_DEFAULT_VERSION) // dot程序版本
	opversion_b := iendecode.Uint64ToBytes(1)                       // 操作版本

	// 打开文件写入
	dop_data_f, err := os.OpenFile(fpath+fname_data, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	defer dop_data_f.Close()
	dop_context_f, err := os.OpenFile(fpath+fname_context, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	defer dop_context_f.Close()

	// 开始写data文件
	_, err = dop_data_f.Write(dotversion_b) // 应用版本
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_data_f.Write(dop.idToByte255(id)) // ID
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_data_f.Write(optime_b) // 时间
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_data_f.Write(opversion_b) // 操作版本
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_data_f.Write(data) // 数据
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}

	// 开始写context文件
	context_liss := []string{}
	context_liss_b, err := iendecode.SliceToBytes("[]string", context_liss)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_context_f.Write(dotversion_b) // 应用版本
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_context_f.Write(optime_b) // 时间
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_context_f.Write(opversion_b) // 操作版本
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_context_f.Write(context_liss_b) // 上下文
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}

	return
}

// 新建一个有一条上下文关系的dot
func (dop *DotsOp) NewDotWithContextT(id string, data []byte, contextid string, context *Context, optime time.Time) (err error) {
	fname, fpath, err := dop.findFilePath(id)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	if len([]byte(contextid)) > DOT_ID_MAX_LENGTH_V1 {
		err = fmt.Errorf("dot: The dot context id length must less than %v: \"%v\"", DOT_ID_MAX_LENGTH_V1, contextid)
		return
	}

	fmt.Println(fpath)

	// 加锁
	dop.dots_lock_lock.Lock()
	if _, have := dop.dots_lock[fname]; have != true {
		dop.dots_lock[fname] = new(sync.RWMutex)
	}
	dop.dots_lock[fname].Lock()
	dop.dots_lock_lock.Unlock()
	// 函数退出的解锁
	defer func() {
		dop.dots_lock_lock.Lock()
		dop.dots_lock[fname].Unlock()
		delete(dop.dots_lock, fname)
		dop.dots_lock_lock.Unlock()
	}()

	fname_data := fname + DOT_FILE_NAME_DATA
	fname_context := fname + DOT_FILE_NAME_CONTEXT
	fname_context_this := fname + DOT_FILE_NAME_CONTEXT + "_" + base.GetSha1Sum(contextid)
	ishave := base.FileExist(fpath + fname_data)
	if ishave == true {
		err = fmt.Errorf("dot: The dot id \"%v\" already have.", id)
		return
	}

	optime_b, _ := optime.MarshalBinary()                           // 操作时间，15的长度
	dotversion_b := iendecode.Uint8ToBytes(DOT_NOW_DEFAULT_VERSION) // dot程序版本
	opversion_b := iendecode.Uint64ToBytes(1)                       // 操作版本

	// 打开文件写入
	dop_data_f, err := os.OpenFile(fpath+fname_data, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	defer dop_data_f.Close()
	dop_context_f, err := os.OpenFile(fpath+fname_context, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	defer dop_context_f.Close()
	dop_context_this_f, err := os.OpenFile(fpath+fname_context_this, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	defer dop_context_this_f.Close()

	// 开始写data文件
	_, err = dop_data_f.Write(dotversion_b) // 应用版本
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_data_f.Write(dop.idToByte255(id)) // ID
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_data_f.Write(optime_b) // 时间
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_data_f.Write(opversion_b) // 操作版本
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_data_f.Write(data) // 数据
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	// 开始写context文件
	context_liss := []string{contextid}
	context_liss_b, err := iendecode.SliceToBytes("[]string", context_liss)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_context_f.Write(dotversion_b) // 应用版本
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_context_f.Write(optime_b) // 时间
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_context_f.Write(opversion_b) // 操作版本
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_context_f.Write(context_liss_b) // 上下文
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	// 开始写context_this文件
	context_b, err := context.MarshalBinary()
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_context_this_f.Write(dotversion_b) // 应用版本
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_context_this_f.Write(dop.idToByte255(contextid)) // Context ID
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_context_this_f.Write(optime_b) // 时间
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_context_this_f.Write(opversion_b) // 操作版本
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = dop_context_this_f.Write(context_b) // 上下文
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}

	return
}

// 完全更新一组上下文
func (dop *DotsOp) UpdateContextT(dotid string, contextid string, context *Context, optime time.Time) (err error) {
	fname, fpath, err := dop.findFilePath(dotid)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	fmt.Println(fpath)

	// 加锁
	dop.dots_lock_lock.Lock()
	if _, have := dop.dots_lock[fname]; have != true {
		dop.dots_lock[fname] = new(sync.RWMutex)
	}
	dop.dots_lock[fname].Lock()
	dop.dots_lock_lock.Unlock()
	// 函数退出的解锁
	defer func() {
		dop.dots_lock_lock.Lock()
		dop.dots_lock[fname].Unlock()
		delete(dop.dots_lock, fname)
		dop.dots_lock_lock.Unlock()
	}()

	// 找有没有文件
	fname_one_context := fname + DOT_FILE_NAME_CONTEXT + "_" + base.GetSha1Sum(contextid)
	ishave := base.FileExist(fpath + fname_one_context)
	if ishave != true {
		err = fmt.Errorf("dot: Can not find the dot \"%v\", or can not find the context \"%v\".", dotid, contextid)
		return
	}

	// 打开文件
	f_con, err := os.OpenFile(fpath+fname_one_context, os.O_RDWR, 0600)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	defer f_con.Close()

	// 构建当前操作时间
	optime_b, _ := optime.MarshalBinary() // 操作时间，15的长度

	var read_n int
	// 获取操作时间
	old_optime_b := make([]byte, 15)
	read_n, err = f_con.ReadAt(old_optime_b, 1+DOT_ID_MAX_LENGTH_V1)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	if read_n != 15 {
		err = fmt.Errorf("dot: not the leght %v")
		return
	}
	old_optime := time.Time{}
	err = old_optime.UnmarshalBinary(old_optime_b)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	// 比较操作时间
	if optime.UnixNano() <= old_optime.UnixNano() {
		err = fmt.Errorf("dot: the opreate lock is wrong for dot %v.", dotid)
		return
	}

	// 获取操作版本
	opversion_b := make([]byte, 8)
	read_n, err = f_con.ReadAt(opversion_b, 1+DOT_ID_MAX_LENGTH_V1+15)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	if read_n != 8 {
		err = fmt.Errorf("dot: net the leght %v")
		return
	}
	opversion := iendecode.BytesToUint64(opversion_b)
	opversion++
	opversion_b = iendecode.Uint64ToBytes(opversion)

	context_b, err := context.MarshalBinary()
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}

	// 准备写入数据
	w_b := make([]byte, 0)
	w_b = append(w_b, optime_b...)
	w_b = append(w_b, opversion_b...)
	w_b = append(w_b, context_b...)
	// 扔掉之前的数据部分
	err = f_con.Truncate(1 + DOT_ID_MAX_LENGTH_V1 + 15 + 8)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	// 重新写
	_, err = f_con.WriteAt(w_b, 1+DOT_ID_MAX_LENGTH_V1)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}

	return
}

// 更新一个数据
func (dop *DotsOp) UpdateDataT(id string, data []byte, optime time.Time) (err error) {
	fname, fpath, err := dop.findFilePath(id)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	fmt.Println(fpath)

	// 加写锁
	dop.dots_lock_lock.Lock()
	if _, have := dop.dots_lock[fname]; have != true {
		dop.dots_lock[fname] = new(sync.RWMutex)
	}
	dop.dots_lock[fname].Lock()
	dop.dots_lock_lock.Unlock()
	// 函数退出的解锁
	defer func() {
		dop.dots_lock_lock.Lock()
		dop.dots_lock[fname].Unlock()
		delete(dop.dots_lock, fname)
		dop.dots_lock_lock.Unlock()
	}()

	// 构建文件名
	fname_data := fname + DOT_FILE_NAME_DATA
	//fname_context := fname + DOT_FILE_NAME_CONTEXT
	ishave := base.FileExist(fpath + fname_data)
	// 如果不存在就返回错误
	if ishave != true {
		err = fmt.Errorf("dot: Can not find the dot \"%v\".", id)
		return
	}

	// 打开数据文件写入
	dop_data_f, err := os.OpenFile(fpath+fname_data, os.O_RDWR, 0600)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	defer dop_data_f.Close()

	// 构建当前操作时间
	optime_b, _ := optime.MarshalBinary() // 操作时间，15的长度

	var read_n int
	// 获取操作时间
	old_optime_b := make([]byte, 15)
	read_n, err = dop_data_f.ReadAt(old_optime_b, 1+DOT_ID_MAX_LENGTH_V1)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	if read_n != 15 {
		err = fmt.Errorf("dot: not the leght %v")
		return
	}
	old_optime := time.Time{}
	err = old_optime.UnmarshalBinary(old_optime_b)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	// 比较操作时间
	if optime.UnixNano() <= old_optime.UnixNano() {
		err = fmt.Errorf("dot: the opreate lock is wrong for dot %v.", id)
		return
	}

	// 获取操作版本
	opversion_b := make([]byte, 8)
	read_n, err = dop_data_f.ReadAt(opversion_b, 1+DOT_ID_MAX_LENGTH_V1+15)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	if read_n != 8 {
		err = fmt.Errorf("dot: net the leght %v")
		return
	}
	opversion := iendecode.BytesToUint64(opversion_b)
	opversion++
	fmt.Println(opversion)
	opversion_b = iendecode.Uint64ToBytes(opversion)

	// 准备写入数据
	w_b := make([]byte, 0)
	w_b = append(w_b, optime_b...)
	w_b = append(w_b, opversion_b...)
	w_b = append(w_b, data...)
	// 扔掉之前的数据部分
	err = dop_data_f.Truncate(1 + DOT_ID_MAX_LENGTH_V1 + 15 + 8)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	// 重新写
	_, err = dop_data_f.WriteAt(w_b, 1+DOT_ID_MAX_LENGTH_V1)

	return
}

// 更新一个上下文中的Down关系
func (dop *DotsOp) UpdateOneDownT(dotid string, contextid string, downname string, varlue string, optime time.Time) (err error) {
	fname, fpath, err := dop.findFilePath(dotid)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	fmt.Println(fpath)

	// 加锁
	dop.dots_lock_lock.Lock()
	if _, have := dop.dots_lock[fname]; have != true {
		dop.dots_lock[fname] = new(sync.RWMutex)
	}
	dop.dots_lock[fname].Lock()
	dop.dots_lock_lock.Unlock()
	// 函数退出的解锁
	defer func() {
		dop.dots_lock_lock.Lock()
		dop.dots_lock[fname].Unlock()
		delete(dop.dots_lock, fname)
		dop.dots_lock_lock.Unlock()
	}()

	// 找有没有文件
	fname_one_context := fname + DOT_FILE_NAME_CONTEXT + "_" + base.GetSha1Sum(contextid)
	ishave := base.FileExist(fpath + fname_one_context)
	if ishave != true {
		err = fmt.Errorf("dot: Can not find the dot \"%v\", or can not find the context \"%v\".", dotid, contextid)
		return
	}

	// 打开文件
	f, err := os.OpenFile(fpath+fname_one_context, os.O_RDWR, 0600)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	defer f.Close()

	// 构建当前操作时间
	optime_b, _ := optime.MarshalBinary() // 操作时间，15的长度

	var read_n int
	// 获取操作时间
	old_optime_b := make([]byte, 15)
	read_n, err = f.ReadAt(old_optime_b, 1+DOT_ID_MAX_LENGTH_V1)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	if read_n != 15 {
		err = fmt.Errorf("dot: not the leght %v")
		return
	}
	old_optime := time.Time{}
	err = old_optime.UnmarshalBinary(old_optime_b)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	// 比较操作时间
	if optime.UnixNano() <= old_optime.UnixNano() {
		err = fmt.Errorf("dot: the opreate lock is wrong for dot %v.", dotid)
		return
	}

	// 获取操作版本
	opversion_b := make([]byte, 8)
	read_n, err = f.ReadAt(opversion_b, 1+DOT_ID_MAX_LENGTH_V1+15)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	if read_n != 8 {
		err = fmt.Errorf("dot: net the leght %v")
		return
	}
	opversion := iendecode.BytesToUint64(opversion_b)
	opversion++
	opversion_b = iendecode.Uint64ToBytes(opversion)

	// 读取这个context的结构体
	con_b, len, err := dop.readAfterWithFile(1+DOT_ID_MAX_LENGTH_V1+15+8, f)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	if len == 0 {
		err = fmt.Errorf("dot: The context have some error.")
		return
	}
	context := &Context{}
	err = context.UnmarshalBinary(con_b)

	// 修改
	context.Down[downname] = varlue
	context_b, err := context.MarshalBinary()
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}

	// 准备写入数据
	w_b := make([]byte, 0)
	w_b = append(w_b, optime_b...)
	w_b = append(w_b, opversion_b...)
	w_b = append(w_b, context_b...)
	// 扔掉之前的数据部分
	err = f.Truncate(1 + DOT_ID_MAX_LENGTH_V1 + 15 + 8)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	// 重新写
	_, err = f.WriteAt(w_b, 1+DOT_ID_MAX_LENGTH_V1)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}

	return
}

// 更新一个上下文的Up关系
func (dop *DotsOp) UpdateOneUpT(dotid string, contextid string, up string, optime time.Time) (err error) {
	fname, fpath, err := dop.findFilePath(dotid)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	fmt.Println(fpath)

	// 加锁
	dop.dots_lock_lock.Lock()
	if _, have := dop.dots_lock[fname]; have != true {
		dop.dots_lock[fname] = new(sync.RWMutex)
	}
	dop.dots_lock[fname].Lock()
	dop.dots_lock_lock.Unlock()
	// 函数退出的解锁
	defer func() {
		dop.dots_lock_lock.Lock()
		dop.dots_lock[fname].Unlock()
		delete(dop.dots_lock, fname)
		dop.dots_lock_lock.Unlock()
	}()

	// 找有没有文件
	fname_one_context := fname + DOT_FILE_NAME_CONTEXT + "_" + base.GetSha1Sum(contextid)
	ishave := base.FileExist(fpath + fname_one_context)
	if ishave != true {
		err = fmt.Errorf("dot: Can not find the dot \"%v\", or can not find the context \"%v\".", dotid, contextid)
		return
	}

	// 打开文件
	f, err := os.OpenFile(fpath+fname_one_context, os.O_RDWR, 0600)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	defer f.Close()

	// 构建当前操作时间
	optime_b, _ := optime.MarshalBinary() // 操作时间，15的长度

	var read_n int
	// 获取操作时间
	old_optime_b := make([]byte, 15)
	read_n, err = f.ReadAt(old_optime_b, 1+DOT_ID_MAX_LENGTH_V1)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	if read_n != 15 {
		err = fmt.Errorf("dot: not the leght %v")
		return
	}
	old_optime := time.Time{}
	err = old_optime.UnmarshalBinary(old_optime_b)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	// 比较操作时间
	if optime.UnixNano() <= old_optime.UnixNano() {
		err = fmt.Errorf("dot: the opreate lock is wrong for dot %v.", dotid)
		return
	}

	// 获取操作版本
	opversion_b := make([]byte, 8)
	read_n, err = f.ReadAt(opversion_b, 1+DOT_ID_MAX_LENGTH_V1+15)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	if read_n != 8 {
		err = fmt.Errorf("dot: net the leght %v")
		return
	}
	opversion := iendecode.BytesToUint64(opversion_b)
	opversion++
	opversion_b = iendecode.Uint64ToBytes(opversion)

	// 读取这个context的结构体
	con_b, len, err := dop.readAfterWithFile(1+DOT_ID_MAX_LENGTH_V1+15+8, f)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	if len == 0 {
		err = fmt.Errorf("dot: The context have some error.")
		return
	}
	context := &Context{}
	err = context.UnmarshalBinary(con_b)

	// 修改
	context.Up = up
	context_b, err := context.MarshalBinary()
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}

	// 准备写入数据
	w_b := make([]byte, 0)
	w_b = append(w_b, optime_b...)
	w_b = append(w_b, opversion_b...)
	w_b = append(w_b, context_b...)
	// 扔掉之前的数据部分
	err = f.Truncate(1 + DOT_ID_MAX_LENGTH_V1 + 15 + 8)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	// 重新写
	_, err = f.WriteAt(w_b, 1+DOT_ID_MAX_LENGTH_V1)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}

	return
}

// 获得Context索引体中的日期和操作版本
func (dop *DotsOp) ReadContextIndexTimeVersionB(dotid string) (time_b []byte, version_b []byte, err error) {
	fname, fpath, err := dop.findFilePath(dotid)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	fmt.Println(fpath)

	// 加读锁
	dop.dots_lock_lock.Lock()
	if _, have := dop.dots_lock[fname]; have != true {
		dop.dots_lock[fname] = new(sync.RWMutex)
	}
	dop.dots_lock[fname].RLock()
	dop.dots_lock_lock.Unlock()
	// 函数退出的解锁
	defer func() {
		dop.dots_lock_lock.Lock()
		dop.dots_lock[fname].RUnlock()
		delete(dop.dots_lock, fname)
		dop.dots_lock_lock.Unlock()
	}()

	// 打开文件
	fname_context := fname + DOT_FILE_NAME_CONTEXT
	ishave := base.FileExist(fpath + fname_context)
	// 如果不存在就返回错误
	if ishave != true {
		err = fmt.Errorf("dot: Can not find the dot \"%v\".", dotid)
		return
	}
	f, err := os.OpenFile(fpath+fname_context, os.O_RDONLY, 0600)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	defer f.Close()

	time_b = make([]byte, 15)
	version_b = make([]byte, 8)

	_, err = f.ReadAt(time_b, 1)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = f.ReadAt(version_b, 1+15)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}

	return
}

// 获得某个Context体中的日期和操作版本
func (dop *DotsOp) ReadContextTimeVersionB(dotid string, contextid string) (time_b []byte, version_b []byte, err error) {
	fname, fpath, err := dop.findFilePath(dotid)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	fmt.Println(fpath)

	// 加读锁
	dop.dots_lock_lock.Lock()
	if _, have := dop.dots_lock[fname]; have != true {
		dop.dots_lock[fname] = new(sync.RWMutex)
	}
	dop.dots_lock[fname].RLock()
	dop.dots_lock_lock.Unlock()
	// 函数退出的解锁
	defer func() {
		dop.dots_lock_lock.Lock()
		dop.dots_lock[fname].RUnlock()
		delete(dop.dots_lock, fname)
		dop.dots_lock_lock.Unlock()
	}()

	// 打开文件
	fname_data := fname + DOT_FILE_NAME_CONTEXT + "_" + base.GetSha1Sum(contextid)
	ishave := base.FileExist(fpath + fname_data)
	// 如果不存在就返回错误
	if ishave != true {
		err = fmt.Errorf("dot: Can not find the dot \"%v\", or can not find the context \"%v\".", dotid, contextid)
		return
	}
	f, err := os.OpenFile(fpath+fname_data, os.O_RDONLY, 0600)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	defer f.Close()

	time_b = make([]byte, 15)
	version_b = make([]byte, 8)

	_, err = f.ReadAt(time_b, 1+DOT_ID_MAX_LENGTH_V1)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = f.ReadAt(version_b, 1+DOT_ID_MAX_LENGTH_V1+15)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}

	return
}

// 获得Data体中的日期和操作版本
func (dop *DotsOp) ReadDataTimeVersionB(dotid string) (time_b []byte, version_b []byte, err error) {
	fname, fpath, err := dop.findFilePath(dotid)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	fmt.Println(fpath)

	// 加读锁
	dop.dots_lock_lock.Lock()
	if _, have := dop.dots_lock[fname]; have != true {
		dop.dots_lock[fname] = new(sync.RWMutex)
	}
	dop.dots_lock[fname].RLock()
	dop.dots_lock_lock.Unlock()
	// 函数退出的解锁
	defer func() {
		dop.dots_lock_lock.Lock()
		dop.dots_lock[fname].RUnlock()
		delete(dop.dots_lock, fname)
		dop.dots_lock_lock.Unlock()
	}()

	// 打开文件
	fname_data := fname + DOT_FILE_NAME_DATA
	ishave := base.FileExist(fpath + fname_data)
	// 如果不存在就返回错误
	if ishave != true {
		err = fmt.Errorf("dot: Can not find the dot \"%v\".", dotid)
		return
	}
	f, err := os.OpenFile(fpath+fname_data, os.O_RDONLY, 0600)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	defer f.Close()

	time_b = make([]byte, 15)
	version_b = make([]byte, 8)

	_, err = f.ReadAt(time_b, 1+DOT_ID_MAX_LENGTH_V1)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	_, err = f.ReadAt(version_b, 1+DOT_ID_MAX_LENGTH_V1+15)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}

	return
}
