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
	"github.com/idcsource/insight00-lib/jconf"
)

// 对某个block进行操作
type BlockOp struct {
	path           string             // block的位置，这两个都要与InitBlock一致
	version        uint8              // block的版本
	deep           uint8              // block的路径深度，这两个都要与InitBlock一致
	dots_lock      map[string]DotLock // 正在操作的dot都会加上相应的锁，map的key为dot的id
	dots_lock_lock *sync.RWMutex      // 避免操作上面的dot锁时有抢占，在对上面的锁修改时也要现锁定
}

// dot的操作锁
type DotLock struct {
	LockTime time.Time
	Lock     *sync.RWMutex
}

// 返回要操作的dot的文件名和路径
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
