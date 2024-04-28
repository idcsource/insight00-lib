# Package webs
一个Web服务器的实现

## 如何使用

### 准备工作

首先需要准备一个配置和一个日志，其中配置的实现见jconf包，日志的实现见logs包。如果你实在不想准备一个日志，那么提供nil也行，那么创建web服务的时候则会建一个运行时日志。

配置中必须包括的信息有：

1. port : Web服务的端口号

2. https : 是否启用https，如果启用，则必须有且必须为true

3. sslcert和sslkey : 启动https所需要的两个证书文件的存放位置

配置中可选的信息有：

1. static : 静态资源文件的总路径，如果不设置则默认为程序所在路径

2. max_routine : 最大并发数，如果不设置则不限制最大并发

### 创建Web实例

通过 NewWeb(conf, log) 可以创建一个Web实例，将会返回*Web供继续配置使用。目前这个服务器实例还无法启动，因为还有后续配置需要进行。

### 后续配置

必须的后续配置如下：

通过InitRouter()方法初始化路由。需要提供FloorInterface接口的执行实例和相关配置文件，程序将生成路由节点树并返回根节点。之后则需要在此基础上配置整个站点的路由，见后面的配置路由详解。

可选的后续配置如下：

1. 通过RegDB()方法注册主数据库，在webs里的默认数据库连接方法是go语言自己提供的database/sql包，所以需要提供*sql.DB来注册。

2. 通过RegMultiDB()方法来注册多个扩展数据库。

3. 通过RegExt()方法来注册扩展，扩展接受所有interface{}类型，使用时需要自己转换格式，如果要注册其他数据方法，也可以在这里注册。

4. 通过RegExecPoint()方法来注册执行点，执行点需要符合ExecPointer接口。

5. 通过ViewPolymer()方法来注册视图聚合，视图聚合需要符合ViewPolymerExecer接口。

6. 通过SetNotFound()方法修改默认的404处理。

7. 通过AddStatic()方法添加静态资源路径，url是浏览器用哪个路径进行访问，path则是相对于前面static的服务器存放位置。例如，static的配置为“/home/web/static/”，而这个路径下有abc和bcd两个路径存放静态文件，那么可以通过AddStatic("st1","abc")和AddStatic("st2","bcd")来注册，浏览器将可以通过http://domain/st1/xxx和http://domain/st2/xxx访问到这两个路径下来的具体静态文件。

### 启动服务

通过Start()就可以启动这个web服务。

### 配置路由详解

需要通过InitRouter()初始化路由之后才能启动web服务，此时当浏览器直接访问http://domain或默认的http://domain/Index时，则由注册的FloorInterface接口执行实例进行执行，并在浏览器上显示返回的页面。

你需要通过InitRouter()所返回的根节点，通过增加门节点、普通节点、空节点逐级构建节点树，或在某个节点树的末尾再注册一个静态文件节点。

#### 普通节点

通过*NodeTree.AddNode(name, mark, floor, config)方法增加一个普通节点。其中name可以用一个对你友好的好任何名字，mark则是相对访问路径，floor则是符合FloorInterface接口的执行实例，config则是jconf的针对这个节点的配置文件。

假设此时你位于根节点，你新增节点的mark输入的是“NewsList”，则这个节点的浏览器访问路径则为：http://domain/NewsList。

添加成为节点后，方法将会返回此节点下的节点树，你可以在此处继续添加下层节点。比如你在此处再添加一个mark是“SoftIntroduce”的节点，则浏览器访问路径则是：http://domain/NewsList/SoftIntroduce。

#### 门节点

通过*NodeTree.AddDoor(name, mark, floordoor, config)方法增加一个门节点。floor则是符合FloorDoorInterface接口的执行实例，其他与普通节点一样。门节点记录了一组平级的普通节点执行实例。

假设此时你位于根节点，你新增门节点的mark输入的是“Products”，而这个门节点记录了“Index”“List”“Detail”三个执行实例，那么浏览器访问路径将分别为：http://domain/Products/Index，http://domain/Products/List，http://domain/Products/Detail。

添加成为节点后，方法将会返回此节点下的节点树，你仍然可以在此处继续添加下层节点，但请确保下层节点的mark不要与门节点中出现的同名。

#### 空节点

通过*NodeTree.AddEmpty(name, mark)方法增加一个空节点。空节点无法独立使用，必须在返回的节点树下增加可用的普通节点、门节点或静态文件节点。

假设此时你位于根节点，你新增空节点的mark输入的是“About”，则浏览器访问路径则为http://domain/About，但这是一个空节点，直接访问将返回404。

#### 静态文件节点

通过*NodeTree.AddStatic(mark, path)方法增加一个静态文件节点。path依然是相对于static配置的相对路径。在此节点下无法再增加新的节点。

### GET风格

本Web服务器使用http://domain/Page1/Page2/:key1=value1/:key2=value2/的形式记录get值，并且通过Runtime.UrlRequest提供给普通节点使用。

### 运行时数据

本Web服务器提供了许多运行时环境数据，具体在const_struct.go中都有注明。

	AllRoutePath string            //整个的RoutePath，也就是除域名外的完整路径
	NowRoutePath []string          //AllRoutePath经过层级路由之后剩余的部分
	RealNode     string            //当前节点的树名，如/node1/node2，如果没有使用节点则此处为空
	MyConfig     *jconf.JsonConf   //当前节点的配置文件，从ConfigTree中获取，如当前节点没有配置文件，则去寻找父节点，直到载入站点的配置文件
	UrlRequest   map[string]string //Url请求的整理，风格为:id=1/:type=notype
	Log          *logs.Logs        // 日志，也就是新建web实例时提供的日志

