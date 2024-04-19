// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ]
// Stephen Fire Meditation Qin [ 火志溟 ] -> stephenfmqin@gmail.com
// This source code is governed by GNU LGPL v3 license

package dot

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/idcsource/insight00-lib/base"
	"github.com/idcsource/insight00-lib/jconf"
)

// 初始化一个块（block）的结构，用来存储dot
// 其中path为保存的路径（需要保证存在），name为存储的名称（会新建），deep为存储目录深度
func InitBlock(path string, name string, deep uint8) (err error) {
	path = base.LocalPath(path)

	path_info, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	if path_info.IsDir() != true {
		err = fmt.Errorf("Dot Block: The \"%v\" not a path", err)
		return
	}

	local_path := path + name + "/"
	err = os.Mkdir(local_path, 0700)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	deployed_file := local_path + DEPLOYED_FILE

	if base.FileExist(deployed_file) == true {
		err = fmt.Errorf("Dot Block: The \"%v\" is already a block", local_path)
		return
	} else {
		// 准备block的配置文件
		b_conf := jconf.NewJsonConf()
		b_conf.AddValueInRoot("version", BLOCK_NOW_DEFAULT_VERSION)
		b_conf.AddValueInRoot("deep", deep)
		var b_conf_j string
		b_conf_j, err = b_conf.OutputJson()
		if err != nil {
			err = fmt.Errorf("Dot Block: %v", err)
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
				err = fmt.Errorf("Dot Block: %v", err)
				break
			}
		}
	}

	return
}

// 启动一个block，将会返回DotsOp。
// 这里的path和name需要填写InitBlock时对应的，也就是这个方法后续将是操作这个Block下的dot。
// 函数将生成一个RUNNING_FILE文件，来确保为正在运行。
func StartBlock(path string, name string) (bop *BlockOp, err error) {
	path = base.LocalPath(path)
	path = path + name + "/"

	// 判断是否为block
	isblock := base.FileExist(path + DEPLOYED_FILE)
	if isblock == false {
		err = fmt.Errorf("Dot Block: This is not a block path: %v", path)
		return
	}
	// 判断是否在运行
	isrunning := base.FileExist(path + RUNNING_FILE)
	if isrunning == true {
		err = fmt.Errorf("Dot Block: The block path is running: %v", path)
		return
	}
	// 获取配置
	b_conf := jconf.NewJsonConf()
	err = b_conf.ReadFile(path + DEPLOYED_FILE)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	b_conf_version, err := b_conf.GetInt64("version")
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	b_conf_deep, err := b_conf.GetInt64("deep")
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	// 生成DotsOp
	bop = &BlockOp{
		path:           path,
		version:        uint8(b_conf_version),
		deep:           uint8(b_conf_deep),
		dots_lock:      make(map[string]*DotLock),
		dots_lock_lock: new(sync.RWMutex),
	}

	// 写入running标记
	f_byte := []byte("1")
	err = ioutil.WriteFile(path+RUNNING_FILE, f_byte, 0600)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}

	return
}

// 关闭block，删除running标记
func (bop *BlockOp) StopBlock() (err error) {
	err = os.Remove(bop.path + RUNNING_FILE)
	return
}

// 备份，需要非running状态下，把整个目录进行打包，需要在linux下，并且有tar和gzip
func BackupBlock(path, name, backupname string) (err error) {
	path = base.LocalPath(path)
	pathname := path + name + "/"

	// 判断是否为block
	isblock := base.FileExist(pathname + DEPLOYED_FILE)
	if isblock == false {
		err = fmt.Errorf("Dot Block: This is not a block path: %v", pathname)
		return
	}
	// 判断是否在运行
	isrunning := base.FileExist(pathname + RUNNING_FILE)
	if isrunning == true {
		err = fmt.Errorf("Dot Block: The block path is running: %v", pathname)
		return
	}

	// 写入running标记
	f_byte := []byte("1")
	err = ioutil.WriteFile(pathname+RUNNING_FILE, f_byte, 0600)
	if err != nil {
		err = fmt.Errorf("Dot Block: %v", err)
		return
	}
	defer os.Remove(pathname + RUNNING_FILE)

	//开始压缩
	env := os.Environ()
	procAttr := &os.ProcAttr{
		Env: env,
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
	}
	// 执行外部压缩命令
	_, err = os.StartProcess("/bin/tar", []string{"tar", "-zcPf", backupname, "--exclude=running", "-C", path, name}, procAttr)
	if err != nil {
		err = fmt.Errorf("Error %v starting process!", err)
	}

	return
}
