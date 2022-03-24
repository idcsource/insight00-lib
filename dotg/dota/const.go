// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package dota

// 用户权限
type UserAuthority uint8

// 操作类型
type OperateType uint8

const (
	DOT_AREA_VERSION       uint8 = 1                       // 目前Dot-Area的版本
	DEPLOYED_FILE                = "dota_deployed"         // 占有文件
	DEFAULT_AREA_BLOCK           = "dota"                  // 默认的提供给Dot-Area使用的block
	DEFAULT_ADMIN_USER           = "insight00"             // 默认的管理员账户，不一定用到
	DEFAULT_ADMIN_PASSWORD       = "insight00"             // 默认的管理员密码，不一定用到
	DEFAULT_BLOCK_DEEP     uint8 = 2                       // 默认的block结构深度
	DEFAULT_RUN_LOG              = "/var/log/dota_run.log" // 默认运行日志位置
	DEFAULT_ERR_LOG              = "/var/log/dota_err.log" // 默认错误日志位置

	USER_AUTHORITY_NO     UserAuthority = iota // 没有权限
	USER_AUTHORITY_ADMIN                       // 管理员权限
	USER_AUTHORITY_NORMAL                      // 普通权限

	OPERATE_TYPE_NO                   OperateType = iota // 操作状态，没有状态
	OPERATE_TYPE_LOGIN                                   // 登陆
	OPERATE_TYPE_KEEPLIVE                                // 续期
	OPERATE_TYPE_PASSWORD                                // 修改密码
	OPERATE_TYPE_NEW_USER                                // 添加用户
	OPERATE_TYPE_DEL_USER                                // 删除用户
	OPERATE_TYPE_NEW_BLOCK                               // 新增Block
	OPERATE_TYPE_DEL_BLOCK                               // 删除Block
	OPERATE_TYPE_NEW_DOT                                 // 新建dot
	OPERATE_TYPE_NEW_DOT_WITH_CONTEXT                    // 新建包含一个上下关系的dot
	OPERATE_TYPE_DEL_DOT                                 // 删除dot
	OPERATE_TYPE_UPDATE_DATA                             // 更新数据
	OPERATE_TYPE_READ_DATA                               // 读取数据
	OPERATE_TYPE_UPDATE_ONE_DOWN                         // 更新一个down
	OPERATE_TYPE_UPDATE_ONE_UP                           // 更新一个up
	OPERATE_TYPE_DEL_ONE_DOWN                            // 删除一个down
	OPERATE_TYPE_ADD_CONTEXT                             // 添加一个context
	OPERATE_TYPE_UPDATE_CONTEXT                          // 完整更新一个context
	OPERATE_TYPE_DEL_CONTEXT                             // 删除一个完整context
	OPERATE_TYPE_READ_CONTEXT                            // 读一个完整的context
	OPERATE_TYPE_READ_ONE_UP                             // 读一个up
	OPERATE_TYPE_READ_ONE_DOWN                           // 读一个down
	OPERATE_TYPE_READ_DATA_TV                            // 读data的time和version
	OPERATE_TYPE_READ_INDEX_TV                           // 读context索引的time和version
	OPERATE_TYPE_READ_CONTEXT_TV                         // 读某个context的time和version
)