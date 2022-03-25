# Package dota

Dot-Area（点域）为多个Dot-Block的集合，可理解为一个普通的数据库。

Dot-Area负责管理一台服务器内的所有Dot-Block，提供通过网络进行访问和管理功能所需要的方法。

本package中包含Dot-Area初始化、Dot-Block的管理、Dot的操作、用户权限管理的相关操作的方法。

Web管理界面、Web Service接口将在其他package中提供。

## dota的初始化过程

1.根据所给位置创建路径，并在其中添加area_deployed文件，写入dot-area的当前版本号，以确定这个是已经被dota占有的。

2.创建默认名为dota的block，深度为2。

3.创建名为allblock的dot，在数据体中用[]string的方式存放area中所有的block名。

4.创建名为alluser的dot，并创建一个alluser的context，up为空，down对应每一个用户。

5.创建默认的用户，dot的id为user_用户名，数据体中记录：40位长密码散列 | uint8的权限类型，1为管理员，2为一般

6.如果新增加用户不是管理员，则创建名为block的context，up为空，down为拥有权限的block名

## dota服务器端配置文件示例

	{
		"server" : {
			"port" : 9999, // 监听端口
			"tls" : true, // 是否tls加密
			"pem_file" : "", // pem文件位置，如果tls为ture
			"key_file" : "", // key文件位置，如果tls为ture
		},
		"log" : {
			"run_log" : "/var/log/dota_run.log", // 运行日志位置，没有则走默认
			"err_log" : "/var/log/dota_err.log", // 错误日志位置，没有则走默认
		}，
		"block":{
			"default_deep" : 3  // block的默认路径深度，如果在创建block的时候不指定，则根据这个默认进行，如果这里没有，则根据常量文件进行，之影响配置文件之后创建的block
		}
	}
