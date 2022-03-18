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
	"github.com/idcsource/insight00-lib/jconf"
)

// 对某个block中的dot进行操作
type DotsOp struct {
	path           string                   // block的位置，这两个都要与InitBlock一致
	version        uint8                    // block的版本
	deep           uint8                    // block的路径深度，这两个都要与InitBlock一致
	dots_lock      map[string]*sync.RWMutex // 正在操作的dot都会加上相应的锁，map的key为dot的id
	dots_lock_lock *sync.RWMutex            // 避免操作上面的dot锁时有抢占，在对上面的锁修改时也要现锁定
}

// 新建一个dot的操作（在block范围内的）
func NewDotsOp(path string, name string) (dop *DotsOp, err error) {
	path = base.LocalPath(path)
	path = path + name + "/"

	isblock := base.FileExist(path + DEPLOYED_FILE)
	if isblock == false {
		err = fmt.Errorf("dot: This is not a block path: %v", path)
		return
	}
	// 获取配置
	b_conf := jconf.NewJsonConf()
	err = b_conf.ReadFile(path + DEPLOYED_FILE)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	b_conf_version, err := b_conf.GetInt64("version")
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	b_conf_deep, err := b_conf.GetInt64("deep")
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}

	dop = &DotsOp{
		path:           path,
		version:        uint8(b_conf_version),
		deep:           uint8(b_conf_deep),
		dots_lock:      make(map[string]*sync.RWMutex),
		dots_lock_lock: new(sync.RWMutex),
	}
	return
}

// 返回要操作的dot的文件名和路径
func (dop *DotsOp) findFilePath(id string) (fname string, fpath string, err error) {
	if len([]byte(id)) > DOT_ID_MAX_LENGTH_V1 {
		err = fmt.Errorf("dot: The dot id length must less than %v: \"%v\"", DOT_ID_MAX_LENGTH_V1, id)
		return
	}
	fname = base.GetSha1Sum(id)
	fpath = dop.path
	for i := 0; i < int(dop.deep); i++ {
		fpath = fpath + string(fname[i]) + "/"
	}

	return
}

func (dop *DotsOp) idToByte255(id string) (b []byte) {
	id_b := []byte(id)

	b = make([]byte, DOT_ID_MAX_LENGTH_V1)
	for i := 0; i < len(id_b); i++ {
		b[i] = id[i]
	}
	return
}

func (dop *DotsOp) byte255ToId(b []byte) (id string) {
	var id_b []byte
	for j := 0; j < DOT_ID_MAX_LENGTH_V1; j++ {
		if b[j] != 0 {
			id_b = append(id_b, b[j])
		}
	}
	id = string(id_b)
	return
}

// 新建一个只有数据的dot
func (dop *DotsOp) NewDot(id string, data []byte) (err error) {
	fname, fpath, err := dop.findFilePath(id)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	fmt.Println(fpath)

	// 加锁
	dop.dots_lock_lock.Lock()
	dop.dots_lock[fname] = new(sync.RWMutex)
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

	optime := time.Now()
	optime_b, _ := optime.MarshalBinary()                           // 操作时间，15的长度
	dotversion_b := iendecode.Uint8ToBytes(DOT_NOW_DEFAULT_VERSION) // dot程序版本
	opversion_b := iendecode.Uint64ToBytes(1)                       // 操作版本

	// 打开文件写入
	dop_data_f, err := os.OpenFile(fpath+fname_data, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	defer dop_data_f.Close()
	dop_context_f, err := os.OpenFile(fpath+fname_context, os.O_WRONLY|os.O_CREATE, 0666)
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
func (dop *DotsOp) NewDotWithContext(id string, data []byte, contextid string, up string, down map[string]string) (err error) {
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
	dop.dots_lock[fname] = new(sync.RWMutex)
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

	optime := time.Now()
	optime_b, _ := optime.MarshalBinary()                           // 操作时间，15的长度
	dotversion_b := iendecode.Uint8ToBytes(DOT_NOW_DEFAULT_VERSION) // dot程序版本
	opversion_b := iendecode.Uint64ToBytes(1)                       // 操作版本

	// 打开文件写入
	dop_data_f, err := os.OpenFile(fpath+fname_data, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	defer dop_data_f.Close()
	dop_context_f, err := os.OpenFile(fpath+fname_context, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	defer dop_context_f.Close()
	dop_context_this_f, err := os.OpenFile(fpath+fname_context_this, os.O_WRONLY|os.O_CREATE, 0666)
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
	context := &Context{
		Up:   up,
		Down: down,
	}
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

// 删除一个dot
func (dop *DotsOp) DelDot(id string) (err error) {
	fname, fpath, err := dop.findFilePath(id)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}

	// 加锁
	dop.dots_lock_lock.Lock()
	dop.dots_lock[fname] = new(sync.RWMutex)
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
	fname_context := fname + DOT_FILE_NAME_CONTEXT
	ishave := base.FileExist(fpath + fname_data)
	// 如果不存在这个文件，就什么都不处理
	if ishave != true {
		//err = fmt.Errorf("dot: Can not find the dot \"%v\".", id)
		return
	}

	// 删除文件
	err = os.Remove(fpath + fname_data)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	err = os.Remove(fpath + fname_context)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}

	return
}

// 更新一个数据
func (dop *DotsOp) UpdateData(id string, data []byte) (err error) {
	fname, fpath, err := dop.findFilePath(id)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	fmt.Println(fpath)

	// 加写锁
	dop.dots_lock_lock.Lock()
	dop.dots_lock[fname] = new(sync.RWMutex)
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
	dop_data_f, err := os.OpenFile(fpath+fname_data, os.O_RDWR, 0666)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	defer dop_data_f.Close()

	// 构建当前操作时间
	optime := time.Now()
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
