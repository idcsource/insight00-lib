// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package dota

const (
	DOT_AREA_VERSION       uint8 = 1                       // 目前Dot-Area的版本
	DEPLOYED_FILE                = "dota_deployed"         // 占有文件
	DEFAULT_AREA_BLOCK           = "dota"                  // 默认的提供给Dot-Area使用的block
	DEFAULT_BLOCK_INDEX          = "allblock"              // 在默认区里存放其他所有存在的block索引名称的dot
	DEFAULT_BLOCK_PREFIX         = "block_"                // block的dot的前缀
	DEFAULT_USER_INDEX           = "alluser"               // 在默认区里存放所有用户名的dot
	DEFAULT_USER_CONTEXT         = "alluser"               // 在用户索引里，保存用户关系的上下文
	DEFAULT_USER_PREFIX          = "user_"                 // 用户dot的前缀
	DEFAULT_ADMIN_USER           = "insight00"             // 默认的管理员账户，不一定用到
	DEFAULT_ADMIN_PASSWORD       = "insight00"             // 默认的管理员密码，不一定用到
	DEFAULT_BLOCK_DEEP     uint8 = 2                       // 默认的block结构深度
	DEFAULT_RUN_LOG              = "/var/log/dota_run.log" // 默认运行日志位置
	DEFAULT_ERR_LOG              = "/var/log/dota_err.log" // 默认错误日志位置
	CLIENT_KEEPLIVE_TIME   int   = 900                     // 客户端续期时间，默认是15分钟，也就是900秒
	SERVER_OUTLOG_TIME     int   = 1800                    // 服务端认为客户端多久没有活动就要重新登录，默认是30分钟，也就是1800秒

	USER_AUTHORITY_NO     uint8 = iota // 没有权限
	USER_AUTHORITY_ADMIN               // 管理员权限
	USER_AUTHORITY_NORMAL              // 普通权限

	// 客户端请求服务端的操作
	OPERATE_TYPE_NO                   uint8 = iota // 操作状态，没有状态
	OPERATE_TYPE_LOGIN                             // 登陆，使用ns包中的To_Login
	OPERATE_TYPE_KEEPLIVE                          // 续期，直接发送ns包中的Login_Base_Info
	OPERATE_TYPE_CHANGE_PASSWORD                   // 修改密码，使用ns包中的Change_Password
	OPERATE_TYPE_NEW_USER                          // 添加用户，使用ns包中的User_PassWd_Power，这个也是用户dot的数据体
	OPERATE_TYPE_USER_ADD_BLOCK                    // 给用户增加一个bock权限，使用ns包中的User_Block
	OPERATE_TYPE_USER_DEL_BLOCK                    // 给用户减一个block权限，使用ns包中的User_Block
	OPERATE_TYPE_DEL_USER                          // 删除用户，直接加用户名
	OPERATE_TYPE_NEW_BLOCK                         // 新增Block，直接加block名
	OPERATE_TYPE_DEL_BLOCK                         // 删除Block，直接加block名
	OPERATE_TYPE_NEW_DOT                           // 新建dot
	OPERATE_TYPE_NEW_DOT_WITH_CONTEXT              // 新建包含一个上下关系的dot
	OPERATE_TYPE_DEL_DOT                           // 删除dot，直接加dot的名字
	OPERATE_TYPE_UPDATE_DATA                       // 更新数据
	OPERATE_TYPE_READ_DATA                         // 读取数据
	OPERATE_TYPE_UPDATE_ONE_DOWN                   // 更新一个down
	OPERATE_TYPE_UPDATE_ONE_UP                     // 更新一个up
	OPERATE_TYPE_DEL_ONE_DOWN                      // 删除一个down
	OPERATE_TYPE_ADD_CONTEXT                       // 添加一个context
	OPERATE_TYPE_UPDATE_CONTEXT                    // 完整更新一个context
	OPERATE_TYPE_DEL_CONTEXT                       // 删除一个完整context
	OPERATE_TYPE_READ_CONTEXT                      // 读一个完整的context
	OPERATE_TYPE_READ_ONE_UP                       // 读一个up
	OPERATE_TYPE_READ_ONE_DOWN                     // 读一个down
	OPERATE_TYPE_READ_DATA_TV                      // 读data的time和version
	OPERATE_TYPE_READ_INDEX_TV                     // 读context索引的time和version
	OPERATE_TYPE_READ_CONTEXT_TV                   // 读某个context的time和version

	// 服务端返回给客户端的状态
	OPERATE_RETURN_NO               uint8 = iota // 操作的返回状态，没有状态
	OPERATE_RETURN_TYPE_NOT_HAVE                 // 请求的操作不存在
	OPERATE_RETURN_TYPE_FORMAT_ERR               // 请求的操作格式错误
	OPERATE_RETURN_LOGIN_OK                      // 成功登陆，后面将跟服务器返回的uuid
	OPERATE_RETURN_PASSWD_NO                     // 用户名或密码错误
	OPERATE_RETURN_LOGIN_NO                      // 没有登陆，客户端应该赶紧重新发起登录
	OPERATE_RETURN_KEEPLIVE_OK                   // 续期成功
	OPERATE_RETURN_ALL_OK                        // 操作都没问题
	OPERATE_RETURN_ALL_OK_WITH_DATA              // 操作都OK，并带有返回数据
	OPERATE_RETURN_ERROR                         // 操作错误，这个必定带有返回数据
)
