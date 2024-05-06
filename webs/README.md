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

2. max_routine : 最大并发数，如果不设置则默认为 CPU数*1000

### 创建Web实例

通过 NewWeb(conf, log) 可以创建一个Web实例，将会返回*Web供继续配置使用。目前这个服务器实例还无法启动，因为还有后续配置需要进行。

### 后续配置

必须的后续配置如下：

通过InitRouter()方法初始化路由。需要提供FloorInterface接口的执行实例和相关配置文件，程序将生成路由节点树并返回根节点。之后则需要在此基础上配置整个站点的路由，见后面的配置路由详解。

可选的后续配置如下：

1. 通过RegDB()方法注册主数据库，在webs里的默认数据库连接方法是go语言自己提供的database/sql包，所以需要提供*sql.DB来注册。

2. 通过RegMultiDB()方法来注册多个扩展数据库。

3. 通过RegExt()方法来注册扩展，扩展接受所有interface{}类型，如果要注册其他数据方法，可以在这里注册。使用时可以用GetExt()取回。

4. 通过RegExecPoint()方法来注册执行点，执行点需要符合ExecPointer接口，可以在任何地方通过ExecPoint()去调用，但只返回错误信息。

5. 通过ViewPolymer()方法来注册视图聚合器，视图聚合器需要符合ViewPolymerExecer接口。

6. 通过SetNotFound()方法修改默认的404处理。

7. 通过AddStatic()方法添加静态资源路径，url是浏览器用哪个路径进行访问，path则是相对于前面static的服务器存放位置。例如，static的配置为“/home/web/static/”，而这个路径下有abc和bcd两个路径存放静态文件，那么可以通过AddStatic("st1","abc")和AddStatic("st2","bcd")来注册，浏览器将可以通过http://domain/st1/xxx 和http://domain/st2/xxx 访问到这两个路径下来的具体静态文件。

### 启动服务

通过Start()就可以启动这个web服务。

### 配置路由详解

需要通过InitRouter()初始化路由之后才能启动web服务，此时当浏览器直接访问http://domain 或默认的http://domain/Index 时，则由注册的FloorInterface接口执行实例进行执行，并在浏览器上显示返回的页面。

你需要通过InitRouter()所返回的根节点，通过增加门节点、普通节点、空节点逐级构建节点树，或在某个节点树的末尾再注册一个静态文件节点。

#### 普通节点

通过*NodeTree.AddNode(name, mark, floor, config)方法增加一个普通节点。其中name可以用一个对你友好的好任何名字，mark则是相对访问路径，floor则是符合FloorInterface接口的执行实例，config则是jconf的针对这个节点的配置文件。

假设此时你位于根节点，你新增节点的mark输入的是“NewsList”，则这个节点的浏览器访问路径则为：http://domain/NewsList 。

添加成为节点后，方法将会返回此节点下的节点树，你可以在此处继续添加下层节点。比如你在此处再添加一个mark是“SoftIntroduce”的节点，则浏览器访问路径则是：http://domain/NewsList/SoftIntroduce 。

#### 门节点

通过*NodeTree.AddDoor(name, mark, floordoor, config)方法增加一个门节点。floor则是符合FloorDoorInterface接口的执行实例，其他与普通节点一样。门节点记录了一组平级的普通节点执行实例。

假设此时你位于根节点，你新增门节点的mark输入的是“Products”，而这个门节点记录了“Index”“List”“Detail”三个执行实例，那么浏览器访问路径将分别为：http://domain/Products/Index ，http://domain/Products/List ，http://domain/Products/Detail 。

添加成为节点后，方法将会返回此节点下的节点树，你仍然可以在此处继续添加下层节点，但请确保下层节点的mark不要与门节点中出现的同名。

#### 空节点

通过*NodeTree.AddEmpty(name, mark)方法增加一个空节点。空节点无法独立使用，必须在返回的节点树下增加可用的普通节点、门节点或静态文件节点。

假设此时你位于根节点，你新增空节点的mark输入的是“About”，则浏览器访问路径则为http://domain/About ，但这是一个空节点，直接访问将返回404。

#### 静态文件节点

通过*NodeTree.AddStatic(mark, path)方法增加一个静态文件节点。path依然是相对于static配置的相对路径。在此节点下无法再增加新的节点。

### GET风格

本Web服务器使用http://domain/Page1/Page2/:key1=value1/:key2=value2/ 的形式记录get值，并且通过Runtime.UrlRequest提供给普通节点使用。

### POST处理

在节点、执行点中，你可以自己从变量r中（也就是*http.Request）自己获取。比如首先执行r.ParseFrom()或r.ParseMultipartForm()，然后再通过r.Form、r.PostFrom、r.FormValue()等方式读取。然后再使用base包里的InputProcessor进行检查和危险字符过滤。

也可以使用本webs包提供的Field字段工具，这个可以直接从*http.Request获取字段，并根据提供的配置文件对字段进行检查和危险字符过滤。配置文件使用jconf包来实现，例子如下：

	{
	"字段名1": "ture, 类型, 显示出来的名字, 具体说明信息, 最小值, 最大值"   # ture为启用，不启用则为false
	}

但请注意，此套Field工具还有不完善的地方，功能并不完整。


### 运行时数据

本Web服务器提供了许多运行时环境数据，具体在const_struct.go中都有注明。

	AllRoutePath string            //整个的RoutePath，也就是除域名外的完整路径
	NowRoutePath []string          //AllRoutePath经过层级路由之后剩余的部分
	RealNode     string            //当前节点的树名，如/node1/node2，如果没有使用节点则此处为空
	MyConfig     *jconf.JsonConf   //当前节点的配置文件
	UrlRequest   map[string]string //Url请求的整理，风格为:id=1/:type=notype
	Log          *logs.Logs        // 日志，也就是新建web实例时提供的日志

### 关于视图聚合

视图聚合是为了方便实现网站页面中重复出现的部件的处理，例如页面的头部和尾部，只需要写一个或几个聚合器就可以在全站共用。

当路由器找到需要执行的普通节点后，执行程序将尝试查看这个节点是否配置了视图聚合，通过访问FloorInterface接口中的ViewPolymer()方法。如果没有配置视图聚合，则执行程序将按照正常的流程执行这个节点，也就是交由FloorInterface接口中的ExecHTTP()来处理后续行为。

如果执行程序发现这个节点配置了视图聚合，则会去尝试执行视图聚合。此时执行程序将会执行FloorInterface接口中的ViewStream()方法，接受返回的数据流以及要求的聚合器名（需要在之前注册过），并将数据流等信息推给聚合器处理，而聚合器也可以要求返回数据流进入另一个聚合器继续聚合。

### 如何去写普通节点、404节点

在源码const_struct.go文件中定义了FloorInterface接口和Floor数据类型，在源码floor.go文件中提供了Floor的原型。通常情况下，你自己的普通节点和404节点应该首先继承Floor，之后再按照需要改写自己的ExecHTTP()、ViewPolymer()或ViewStream()方法。

### 关于普通节点和门节点的复用

对于一个站点来说，不可避免会出现某些路径下的功能类似或一样的情况。比如不同类型的新闻，展现方式和功能都类似。对于此类情况，你可以使用同一个普通节点或门节点的代码。在添加普通节点或门节点的时候，会被要求提供jconf的配置文件，并通过运行时数据提供给当前运行的节点，你可以用配置文件来告诉代码自己在当前节点下的运行方式。

### 关于自动跳转节点

本Webs服务器还提供了一个MoveToFloor的特殊普通节点，其实现了FloorInterface接口，可以直接使用，并作为普通节点注册进路由中，但注册时需要提供一个跳转节点的路径，例如：

	route_tree.AddDoor("自动跳转", "jump", &webs.MoveToFloor{Url: "/Admin/login"}, nil)
	
程序将在被执行到此节点时，直接给浏览器发送跳转的303指令。

这个功能在通常情况下是没用的。