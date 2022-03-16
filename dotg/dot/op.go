// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package dot

import (
	"fmt"
	"sync"

	"github.com/idcsource/insight00-lib/base"
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
		path:      path,
		version:   uint8(b_conf_version),
		deep:      uint8(b_conf_deep),
		dots_lock: make(map[string]*sync.RWMutex),
	}
	return
}

// 返回要操作的dot的文件名和路径
func (dop *DotsOp) findFilePath(id string) (fname string, fpath string) {
	fpath = dop.path
	for i := 0; i < int(dop.deep); i++ {
		fpath = fpath + string(id[i]) + "/"
	}
	fname = base.GetSha1Sum(id)
	return
}

// 新建一个只有数据的dot
func (dop *DotsOp) NewDot(id string, data []byte) (err error) {
	fname, fpath := dop.findFilePath(id)
	fname_data := fname + DOT_FILE_NAME_DATA
	//fname_context := fname + DOT_FILE_NAME_CONTEXT
	ishave := base.FileExist(fpath + fname_data)
	if ishave == true {
		err = fmt.Errorf("dot: The dot id \"%v\" already have.", id)
		return
	}
	return
}

// 新建一个有一条上下文关系的dot
func (dop *DotsOp) NewDotWithContext(id string, data []byte, context string, up string, down map[string]string) (err error) {
	return
}

// 删除一个dot
func (dop *DotsOp) DelDot(id string) (err error) {
	return
}
