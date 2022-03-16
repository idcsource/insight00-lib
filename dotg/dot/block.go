// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package dot

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/idcsource/insight00-lib/base"
	"github.com/idcsource/insight00-lib/jconf"
)

// 初始化一个块（block）的结构，用来存储dot
func InitBlock(path string, name string, deep uint8) (err error) {
	path = base.LocalPath(path)

	path_info, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	if path_info.IsDir() != true {
		err = fmt.Errorf("dot: The \"%v\" not a path", err)
		return
	}

	local_path := path + name + "/"
	err = os.Mkdir(local_path, 0700)
	if err != nil {
		err = fmt.Errorf("dot: %v", err)
		return
	}
	deployed_file := local_path + DEPLOYED_FILE

	if base.FileExist(deployed_file) == true {
		err = fmt.Errorf("dot: The \"%v\" is already a block", local_path)
		return
	} else {
		// 准备block的配置文件
		b_conf := jconf.NewJsonConf()
		b_conf.AddValueInRoot("version", BLOCK_NOW_DEFAULT_VERSION)
		b_conf.AddValueInRoot("deep", deep)
		var b_conf_j string
		b_conf_j, err = b_conf.OutputJson()
		if err != nil {
			err = fmt.Errorf("dot: %v", err)
			return
		}
		f_byte := []byte(b_conf_j)
		ioutil.WriteFile(deployed_file, f_byte, 0600)
	}

	l_path_name := []string{"a", "b", "c", "d", "e", "f", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	l_path := make([]string, 0)

	for i := 0; i < int(deep); i++ {
		if len(l_path) == 0 {
			for _, v := range l_path_name {
				l_path = append(l_path, local_path+v+"/")
			}
		} else {
			ll_path := make([]string, 0)
			for _, v := range l_path {
				for _, v2 := range l_path_name {
					ll_path = append(ll_path, v+v2+"/")
				}
			}
			l_path = append(l_path, ll_path...)
		}
	}
	for _, v := range l_path {
		if base.FileExist(v) == true {
			continue
		} else {
			err = os.Mkdir(v, 0700)
			if err != nil {
				err = fmt.Errorf("dot: %v", err)
				break
			}
		}
	}

	return
}
