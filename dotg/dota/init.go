// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package dota

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/idcsource/insight00-lib/base"
	"github.com/idcsource/insight00-lib/dotg/dot"
	"github.com/idcsource/insight00-lib/iendecode"
)

// 初始化一个Dot-Area，在制定的目录path下创建area的完整结构，并创建默认的block并创建默认的管理员用户名和密码
func InitArea(path string) (err error) {
	// 查看路径情况
	path = base.LocalPath(path)
	if base.FileExist(path) != true {
		err = os.Mkdir(path, 0700)
		if err != nil {
			err = fmt.Errorf("dota: %v", err)
			return
		}
	} else {
		var path_info fs.FileInfo
		path_info, err = os.Stat(path)
		if err != nil {
			err = fmt.Errorf("dota: %v", err)
			return
		}
		if path_info.IsDir() != true {
			err = fmt.Errorf("dota: The \"%v\" not a path.", err)
			return
		}
	}
	// 写站位文件
	if base.FileExist(path+DEPLOYED_FILE) == true {
		err = fmt.Errorf("dota: The \"%v\" is already a Dot-Area", path)
		return
	}
	area_version_b := iendecode.Uint8ToBytes(DOT_AREA_VERSION)
	err = ioutil.WriteFile(path+DEPLOYED_FILE, area_version_b, 0600)
	if err != nil {
		err = fmt.Errorf("dota: %v", err)
		return
	}

	// 创建默认block
	err = dot.InitBlock(path, DEFAULT_AREA_BLOCK, DEFAULT_BLOCK_DEEP)
	if err != nil {
		err = fmt.Errorf("dota: %v", err)
		return
	}
	// 准备对新的block操作
	op, err := dot.NewDotsOp(path, DEFAULT_AREA_BLOCK)
	if err != nil {
		err = fmt.Errorf("dota: %v", err)
		return
	}
	// 准备block索引的dot
	block_index_data := []string{DEFAULT_AREA_BLOCK}
	block_index_data_b := iendecode.SliceStringToBytes(block_index_data)
	err = op.NewDot(DEFAULT_BLOCK_INDEX, block_index_data_b)
	if err != nil {
		err = fmt.Errorf("dota: %v", err)
		return
	}
	// 准备admin索引的dot
	admin_context := dot.NewContext()
	admin_context.Down[DEFAULT_ADMIN_USER] = DEFAULT_ADMIN_USER
	err = op.NewDotWithContext(DEFAULT_ADMIN_INDEX, []byte(DEFAULT_ADMIN_USER), DEFAULT_ADMIN_CONTEXT, admin_context)
	if err != nil {
		err = fmt.Errorf("dota: %v", err)
		return
	}
	// 准备默认的admin
	admin_data := &Admin_PassWd_Power{
		Password:  base.GetSha1Sum(DEFAULT_ADMIN_PASSWORD),
		PowerType: USER_AUTHORITY_ADMIN,
	}
	admin_data_b, err := admin_data.MarshalBinary()
	if err != nil {
		err = fmt.Errorf("dota: %v", err)
		return
	}
	err = op.NewDot(DEFAULT_ADMIN_PREFIX+DEFAULT_ADMIN_USER, admin_data_b)
	if err != nil {
		err = fmt.Errorf("dota: %v", err)
		return
	}

	return
}
