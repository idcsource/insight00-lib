# Package jconf

一个使用JSON风格语言的配置文件管理。

使用NewJsonConf()方法新建一个实例，之后就可以使用ReadFile()或ReadString()读取JSON配置文件。

JConf使用了节点的方式将具体的配置文件项进行检出，节点可以通过“node1>node2>node3”访问到，返回的类型具体看各方法。

但目前节点不能跨过JSON的数组，JSON数组只能使用GetValue()整体返回，需要自己再进行类型断言。或者也可以使用GetStruct()方法，将数据装入结构体。

AddValue()和DelValue()有针对根节点操作的独立函数，因为不想这两个方法的node参数上出现无意义的nil或""。

SetValue()可以对节点的内容进行更改，并且可以放入任何数据。如果新增加的数据是合法的结构体或map，那么JSON的节点也会跟着调整。AddValue()与AddValueInRoot()也会如此。

GetNode()方法会整体检出节点下的配置文件，并形成新的根节点和JsonConf实例。但请注意，对新JsonConf的所有删减修改操作，都会反映到原有配置文件中。

OutPutJson()将会返回格式化的JSON字符串。
