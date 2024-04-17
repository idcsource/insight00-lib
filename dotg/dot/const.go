// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ]
// Stephen Fire Meditation Qin [ 火志溟 ] -> stephenfmqin@gmail.com
// This source code is governed by GNU LGPL v3 license

package dot

const (
	DOT_NOW_DEFAULT_VERSION uint8 = 2 // 当前dot默认版本号，涉及到后续升级时使用

	BLOCK_NOW_DEFAULT_VERSION uint8 = 1 // 当前block默认版本号，涉及到后续升级时使用

	DOT_FILE_NAME_DATA = "_data"

	DOT_FILE_NAME_CONTEXT = "_context"

	DEPLOYED_FILE = "deployed"

	RUNNING_FILE = "running" // 标记正在运行的文件

	DOT_ID_MAX_LENGTH_V2 = 255 // dot版本号为1与2的id最大长度

	DOT_CONTENT_MAX_IN_DATA_V2 = 255 // dot版本号为2的content内连配置数据的最大长度
)

// dot上下文关系总索引的状态位
const (
	DOT_CONTENT_INDEX_NOTHING uint8 = 1 // 空位不用
	DOT_CONTENT_INDEX_DEL               // 标记删除
)

// dot上下文关系DOWN索引的状态位
const (
	DOT_CONTENT_UP_DOWN_INDEX_NOTHING       uint8 = 1 // 空位不用
	DOT_CONTENT_UP_DOWN_INDEX_DEL                     // 标记删除，UP关系不用
	DOT_CONTENT_UP_DOWN_INDEX_INDATA                  // 数据在里面（不足255bit）
	DOT_CONTENT_UP_DOWN_INDEX_OUTDATA                 // 数据在外面（超过255bit）
	DOT_CONTENT_UP_DOWN_INDEX_OUTDATA_NODEL           // 不需要的外部数据文件但还没有删
)
