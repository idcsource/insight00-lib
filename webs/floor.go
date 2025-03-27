// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package webs

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/idcsource/insight00-lib/base"
)

// 函数负责对控制器进行初始化。
// 根据Runtime中的NowRoutePath，使用strings.Split函数根据“/”分割并作为Runtime中的UrlRequest（已清除可能的为空的字符串）
func (f *Floor) InitHTTP(w http.ResponseWriter, r *http.Request, b *Web, rt Runtime) {
	f.W = w
	f.R = r
	f.Rt = rt
	f.B = b
}

func (f *Floor) ExecHTTP() {

}

func (f *Floor) ViewPolymer() (switchs PolymerSwitch) {
	switchs = POLYMER_NO
	return
}

// order是下一步需要去执行的视图聚合的名称，这个在前期应该被注册过
func (f *Floor) ViewStream() (stream string, order string, data interface{}) {
	return
}

// 无法找到页面的系统内默认处理手段
type NotFoundFloor struct {
	Floor
}

func (n *NotFoundFloor) ExecHTTP() {
	n.W.WriteHeader(404)
	fmt.Fprint(n.W, "404 Page Not Found")
	return
}

// 静态文件的系统内默认处理手段
type StaticFileFloor struct {
	Floor
	path   string
	candir bool // 是否允许列出目录
}

func (f *StaticFileFloor) ExecHTTP() {

	thefile := strings.Join(f.Rt.NowRoutePath, "/")
	thefile = f.B.static + base.DirMustEnd(f.path) + thefile

	finfo, err := os.Lstat(thefile)
	if err != nil {
		f.B.toNotFoundHttp(f.W, f.R, f.Rt) //找不到404
	} else if finfo.IsDir() == true && f.candir == true {
		http.ServeFile(f.W, f.R, thefile) // 是目录但允许列目录
	} else if finfo.IsDir() == true && f.candir == false {
		f.B.toNotFoundHttp(f.W, f.R, f.Rt) // 是目录但不允许列目录
	} else {
		http.ServeFile(f.W, f.R, thefile) // 不是目录
	}
}

// 空节点的处理手段
type EmptyFloor struct {
	Floor
}

func (f *EmptyFloor) ExecHTTP() {
	f.B.toNotFoundHttp(f.W, f.R, f.Rt)
}

// 自动跳转到地址的节点处理手段
type MoveToFloor struct {
	Floor
	Url string
}

func (f *MoveToFloor) ExecHTTP() {
	http.Redirect(f.W, f.R, f.Url, 303)
}
