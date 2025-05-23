// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

// 一个HTTP服务器的实现
package webs

import (
	"database/sql"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/idcsource/insight00-lib/logs"
)

const (
	// 默认最大并发
	MAX_ROUTINE_RATIO = 10
	// 默认页面锁清理循环时间，秒
	DEFAULT_CLEAN_PAGE_LOCK_MAIN = 10
	// 默认的页面锁超时时间，秒
	DEFAULT_PAGE_LOCK_OUTTIME = 15
	// 默认的页面锁延迟，Microsecond，微秒
	DEFAULT_PAGE_LOCK_DELAY = 10
)

const (
	//节点的类型
	NODE_IS_NO     = iota // 没有定义
	NODE_IS_ROOT          // 根结点
	NODE_IS_DOOR          // 入口Door
	NODE_IS_NOMAL         // 普通
	NODE_IS_STATIC        // 静态
	NODE_IS_EMPTY         // 空
)

type PolymerSwitch uint8 // The Polymer Switch

const (
	POLYMER_NO     PolymerSwitch = iota // Not to make polymer
	POLYMER_TYPE_1                      // Use type 1 to make polymer
	POLYMER_TYPE_2                      // Use type 2 to make polymer
	POLYMER_TYPE_3                      // Use type 3 to make polymer
	POLYMER_TYPE_4                      // Use type 4 to make polymer
	POLYMER_TYPE_5                      // Use type 5 to make polymer
)

// Web的数据结构
type Web struct {
	local          string                       // 本地路径
	static         string                       // 静态资源路径
	config         Configer                     // 自身的配置文件
	DB             *sql.DB                      // 主数据库连接，使用Go语言自己提供的方法
	MultiDB        map[string]*sql.DB           // 扩展多数据库准备，使用Go语言自己提供的方法
	ext            map[string]interface{}       // Extension扩展数据（功能）
	execpoint      map[string]ExecPointer       // 执行点
	viewpolymer    map[string]ViewPolymerExecer // view polymer's interface
	router         *Router                      // 路由器
	log            logs.Logser                  // 运行日志
	visit_log      bool                         // 是否开启访问日志
	page_lock_main *sync.Mutex                  // 页面锁总控
	page_lock      map[string]*PageLock         // 页面锁
	max_routine    chan bool                    // 最大并发
}

// 页面锁
type PageLock struct {
	Lock *sync.Mutex
	Time time.Time
}

// 路由器基本类型
type Router struct {
	node_tree    *NodeTree                 // 节点树
	not_found    FloorInterface            // 404路由
	router_ok    bool                      // 其实就是看是否已经设定了NodeTree的根节点
	static_route map[string]*regexp.Regexp // 静态路由
}

// 节点树基本数据类型
type NodeTree struct {
	name        string               // 节点的名称
	mark        string               // 用来做路由的，也就是未来显示在连接上的地址
	config      Configer             // 节点配置文件
	if_children bool                 // 是否有下层
	node_type   int                  // 类型，首页、普通页、入口Door，NODE_IS_*
	floor       FloorInterface       // 控制器
	door        FloorDoorInterface   // 门入口
	children    map[string]*NodeTree // 下层的信息，map的键为Mark
}

// 运行时数据结构
type Runtime struct {
	AllRoutePath string            //整个的RoutePath，也就是除域名外的完整路径
	NowRoutePath []string          //AllRoutePath经过层级路由之后剩余的部分
	RealNode     string            //当前节点的树名，如/node1/node2，如果没有使用节点则此处为空
	WebConfig    Configer          //Web站点的总配置文件
	MyConfig     Configer          //当前节点的配置文件
	UrlRequest   map[string]string //Url请求的整理，风格为:id=1/:type=notype
	Log          logs.Logser       // 日志
}

// FloorInterface 此为控制器接口的定义
type FloorInterface interface {
	InitHTTP(w http.ResponseWriter, r *http.Request, b *Web, rt Runtime)
	ExecHTTP()
	ViewPolymer() (switchs PolymerSwitch)
	ViewStream() (stream string, order string, data interface{})
}

// 控制器原型的数据类型
type Floor struct {
	W  http.ResponseWriter
	R  *http.Request
	Rt Runtime
	B  *Web
}

// FloorDoor的接口和数据类型
type FloorDoor map[string]FloorInterface

type FloorDoorInterface interface {
	FloorList() FloorDoor
}

// 执行点的接口定义
type ExecPointer interface {
	ExecPoint(w http.ResponseWriter, r *http.Request, b *Web, rt Runtime) (err error)
}

// View Polymer's Execer
type ViewPolymerExecer interface {
	Exec(switchs PolymerSwitch, rt Runtime, stream string, data interface{}) (newstream string, newswitchs PolymerSwitch, neworder string, newdata interface{})
}

// 配置接口，支持jconf和yconf两个库，也就是支持JSON或YAML风格的配置文件
type Configer interface {
	AddValue(node string, name string, value interface{}) (err error)
	AddValueInRoot(name string, value interface{}) (err error)
	DelValue(node string, name string) (err error)
	DelValueInRoot(name string) (err error)
	GetArray(node string) (ar []interface{}, err error)
	GetBool(node string) (b bool, err error)
	GetEnum(node string) (em []string, err error)
	GetFloat64(node string) (f64 float64, err error)
	GetInt64(node string) (i64 int64, err error)
	// GetNode(node string) (newjconf Configer, err error)
	GetString(node string) (str string, err error)
	GetStruct(node string, v interface{}) (err error)
	GetValue(node string) (oneNodeVal interface{}, err error)
	// MarshalBinary() (data []byte, err error)
	Println()
	ReadFile(fname string) (err error)
	ReadString(yamlstream string) (err error)
	SetValue(node string, value interface{}) (err error)
	// UnmarshalBinary(data []byte) (err error)
}
