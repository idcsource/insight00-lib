// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package webs

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/idcsource/insight00-lib/base"
	"github.com/idcsource/insight00-lib/logs"
)

// 创建一个Web，db数据库和log日志可以为nil
func NewWeb(config Configer, log logs.Logser) (web *Web) {
	if log == nil {
		log, _ = logs.NewRunLoger(100)
	}
	web = &Web{
		local:          base.LocalPath(""),
		config:         config,
		MultiDB:        make(map[string]*sql.DB),
		ext:            make(map[string]interface{}),
		execpoint:      make(map[string]ExecPointer),
		viewpolymer:    make(map[string]ViewPolymerExecer),
		log:            log,
		visit_log:      true,
		page_lock_main: new(sync.Mutex),
		page_lock:      make(map[string]*PageLock),
		router:         newRouter(),
	}
	// 检查静态资源地址是不是有
	static, err := web.config.GetString("static")
	if err != nil {
		static = web.local
	} else {
		static = base.LocalPath(static)
		static = base.DirMustEnd(static)
	}
	web.static = static
	// 检查是否配备了访问日志
	visit_log, err := web.config.GetBool("visitlog")
	if err == nil {
		web.visit_log = visit_log
	}

	// 准备最大并发
	var max int64
	var ok1 error
	max, ok1 = web.config.GetInt64("max_routine")
	if ok1 != nil {
		max = int64(runtime.NumCPU()) * MAX_ROUTINE_RATIO
	}
	web.max_routine = make(chan bool, max)

	return
}

// 获取本地路径
func (web *Web) GetLocalPath() (path string) {
	return web.local
}

// 获取静态文件路径
func (web *Web) GetStaticPath() (path string) {
	return web.static
}

// 注册主数据库
func (web *Web) RegDB(database *sql.DB) {
	web.DB = database
	return
}

// 注册扩展数据库
func (web *Web) RegMultiDB(name string, db *sql.DB) {
	web.MultiDB[name] = db
}

// 注册扩展
func (web *Web) RegExt(name string, ext interface{}) {
	web.ext[name] = ext
}

// 获取扩展
func (web *Web) GetExt(name string) (ext interface{}, err error) {
	_, find := web.ext[name]
	if find == false {
		err = fmt.Errorf("webs[Web]GetExt: The Extend %v not registered.", name)
		return
	}
	return web.ext[name], nil
}

// Reg View Polymer Execer
func (web *Web) ViewPolymer(name string, execer ViewPolymerExecer) {
	web.viewpolymer[name] = execer
}

// 注册执行点
func (web *Web) RegExecPoint(name string, point ExecPointer) {
	web.execpoint[name] = point
}

// 执行执行点
func (web *Web) ExecPoint(name string, w http.ResponseWriter, r *http.Request, b *Web, rt Runtime) (err error) {
	_, find := web.execpoint[name]
	if find == false {
		return fmt.Errorf("Can not found the Exec Point.")
	}
	return web.execpoint[name].ExecPoint(w, r, b, rt)
}

// 创建路由，设置根节点，并返回根结点，之后所有的对节点的添加操作均是*NodeTree提供的方法
func (web *Web) InitRouter(f FloorInterface, config Configer) (root *NodeTree) {
	return web.router.buildRouter(f, config)
}

// 创建静态地址,path必须是相对于静态地址(static)的地址（不再提供的功能）
// func (web *Web) AddStatic(url, path string) {
// 	path = base.AbsolutePath(path, web.static)
// 	web.router.addStatic(url, path)
// }

// 修改默认的404处理
func (web *Web) SetNotFound(f FloorInterface) {
	web.router.not_found = f
}

func (web *Web) Start() (err error) {
	// 如果没有初始化路由
	if web.router.router_ok == false {
		err = fmt.Errorf("webs[Web]Start: The Router not initialization.")
		web.log.WriteLog(err.Error())
		return
	}

	/* 检查一堆配置文件是否有 */

	// 检查端口是否有
	port, err := web.config.GetString("port")
	if err != nil {
		err = fmt.Errorf("webs[Web]Start: The config port not be set.")
		web.log.WriteLog(err.Error())
		return
	}
	// 检查是http还是https
	var ifHttps bool
	ifHttps, err = web.config.GetBool("https")
	if err != nil {
		ifHttps = false
		err = nil
	}
	var thecert, thekey string
	if ifHttps == true {
		var e2, e3 error
		thecert, e2 = web.config.GetString("sslcert")
		thekey, e3 = web.config.GetString("sslkey")
		if e2 != nil || e3 != nil {
			err = fmt.Errorf("webs[Web]Start: The SSL cert or key not be set !")
			web.log.WriteLog(err.Error())
			return
		}
	}

	/* 启动HTTP服务 */
	port = ":" + port

	go func() {
		if ifHttps == true {
			err = http.ListenAndServeTLS(port, base.LocalFile(thecert), base.LocalFile(thekey), web)
		} else {
			err = http.ListenAndServe(port, web)
		}
		if err != nil {
			err = fmt.Errorf("webs[Web]Start: Can not start the web server: %v", err)
			web.log.WriteLog(err.Error())
			return
		}
	}()

	go web.cleanPageLock() // 页面锁清理子进程

	return
}

// 页面锁清理
func (web *Web) cleanPageLock() {
	for {
		time.Sleep(DEFAULT_CLEAN_PAGE_LOCK_MAIN * time.Second) // 先休息10秒钟

		web.page_lock_main.Lock()

		timenow := time.Now()

		for k, v := range web.page_lock {
			d := timenow.Sub(v.Time).Seconds()
			if d >= DEFAULT_PAGE_LOCK_OUTTIME {
				delete(web.page_lock, k)
			}
		}

		web.page_lock_main.Unlock()
	}
}

// HTTP的路由，提供给"net/http"包使用
func (web *Web) ServeHTTP(httpw http.ResponseWriter, httpr *http.Request) {
	//对进程数的控制
	web.max_routine <- true
	defer func() {
		<-web.max_routine
	}()

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic: ", err)
		}
	}()

	//要运行的Floor
	var runfloor FloorInterface
	//将获得的URL用斜线拆分成[]string
	urla, parameter := base.SplitUrl(httpr.URL.Path)
	//准备基本的RunTime
	rt := Runtime{
		AllRoutePath: httpr.URL.Path,
		NowRoutePath: urla,
		UrlRequest:   parameter,
		WebConfig:    web.config,
		Log:          web.log,
	}

	// 这里有个锁实现，避免并发太集中
	var thelock *PageLock
	var have bool
	web.page_lock_main.Lock()
	thelock, have = web.page_lock[rt.AllRoutePath]
	if have == false {
		web.page_lock[rt.AllRoutePath] = &PageLock{Lock: new(sync.Mutex), Time: time.Now()}
		thelock = web.page_lock[rt.AllRoutePath]
	}
	web.page_lock_main.Unlock()
	thelock.Lock.Lock()
	go func() {
		time.Sleep(DEFAULT_PAGE_LOCK_DELAY * time.Microsecond)
		thelock.Time = time.Now()
		thelock.Lock.Unlock()
	}()

	//静态路由(不再提供的功能)
	// static, have := web.router.getStatic(httpr.URL.Path)
	// if have == true {
	// 	_, err := os.Stat(static)
	// 	if err != nil {
	// 		web.toNotFoundHttp(httpw, httpr, rt)
	// 	} else {
	// 		finfo, err := os.Lstat(static)
	// 		if err != nil {
	// 			web.toNotFoundHttp(httpw, httpr, rt)
	// 		} else if finfo.IsDir() {
	// 			web.toNotFoundHttp(httpw, httpr, rt)
	// 		} else {
	// 			http.ServeFile(httpw, httpr, static)
	// 		}
	// 	}
	// 	return
	// }

	// 如果为0,则处理首页，直接取出NodeTree的根节点
	if len(urla) == 0 {
		rt.RealNode = ""
		runfloor = web.router.node_tree.floor
	} else {
		runfloor, rt = web.router.getRunFloor(rt)
	}

	//开始执行
	runfloor.InitHTTP(httpw, httpr, web, rt)
	switchs := runfloor.ViewPolymer()
	if switchs == POLYMER_NO {
		runfloor.ExecHTTP()
	} else {
		// the view polymer
		var stream, order string
		var data interface{}
		stream, order, data = runfloor.ViewStream()
		if order == "" {
			fmt.Fprint(httpw, stream)
			return
		}
		for {
			oneexec, have := web.viewpolymer[order]
			if have == false {
				fmt.Fprint(httpw, "The ViewPolymer set is wrong, cannot find %v.", order)
				return
			}
			stream, switchs, order, data = oneexec.Exec(switchs, rt, stream, data)
			if switchs == POLYMER_NO {
				break
			}
		}
		fmt.Fprint(httpw, stream)
	}
	return
}

// 去执行NotFound，不要直接调用这个方法
func (web *Web) toNotFoundHttp(w http.ResponseWriter, r *http.Request, rt Runtime) {
	runfloor := web.router.not_found
	runfloor.InitHTTP(w, r, web, rt)
	runfloor.ExecHTTP()
	return
}

// 获取IP的内部函数
func (web *Web) getIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	return "", errors.New("no valid ip found")
}

// 执行记录访问日志
func (web *Web) ToVisitLog(r *http.Request, rt Runtime, stat string) {
	if web.visit_log == false {
		return
	}
	ip, err := web.getIP(r)
	if err != nil {
		ip = "0.0.0.0"
	}
	web.log.WriteLog(ip + " : " + r.Proto + " : " + r.Host + " : " + rt.AllRoutePath + " : " + r.UserAgent() + " : " + stat)
}
