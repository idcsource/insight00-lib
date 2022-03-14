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

// order is the View Polymer Execer's name witch will exec next step.
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
	path string
}

func (f *StaticFileFloor) ExecHTTP() {

	thefile := strings.Join(f.Rt.NowRoutePath, "/")
	thefile = f.B.static + thefile

	_, err := os.Stat(thefile)
	if err != nil {
		f.B.toNotFoundHttp(f.W, f.R, f.Rt)
	} else {
		http.ServeFile(f.W, f.R, thefile)
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
