// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package dot

import (
	"fmt"
	"sync"

	"github.com/idcsource/insight00-lib/base"
)

// 对某个block中的dot进行操作
type DotsOp struct {
	path           string                   // block的位置，这两个都要与InitBlock一致
	deep           uint8                    // block的路径深度，这两个都要与InitBlock一致
	dots_lock      map[string]*sync.RWMutex // 正在操作的dot都会加上相应的锁，map的key为dot的id
	dots_lock_lock *sync.RWMutex            // 避免操作上面的dot锁时有抢占，在对上面的锁修改时也要现锁定
}

func NewDotsOp(path string, name string, deep uint8) (dop *DotsOp, err error) {
	path = base.LocalPath(path)
	path = path + name + "/"

	isblock := base.FileExist(path + DEPLOYED_FILE)
	if isblock == false {
		err = fmt.Errorf("dot: This is not a block path: %v", path)
		return
	}

	dop = &DotsOp{
		path: path,
		deep: deep,
	}
	return
}

func (dop *DotsOp) findFilePath(id string) (fname string, fpath string) {
	fpath = dop.path
	for i := 0; i < int(dop.deep); i++ {
		fpath = fpath + string(id[i]) + "/"
	}
	fname = base.GetSha1Sum(id)
	return
}

func (dop *DotsOp) NewDot(id string, data []byte) (err error) {
	fname, fpath := dop.findFilePath(id)
	fname_data := fname + DOT_FILE_NAME_DATA
	fname_context := fname + DOT_FILE_NAME_CONTEXT
	ishave := base.FileExist(fpath + fname_data)
	if ishave == true {
		err = fmt.Errorf("dot: The dot id \"%v\" already have.", id)
		return
	}
	return
}

func (dop *DotsOp) NewDotWithContext(id string, data []byte, context string, up string, down map[string]string) (err error) {
	return
}
