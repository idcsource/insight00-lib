// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

// 一个最小数据存储单位的实现
package dot

import (
	"time"
)

type Dot struct {
	Version uint8  // 涉及所操作dot的版本号，现在已经为V2
	Id      string // dot的id，固定在255个字节长度的字符串，存储的时候，后面会补0字节
	//	Time     time.Time           // 这次修改dot的时间点，作为一致性检查使用，V2废弃
	Edition  uint64              // 修改版本号，新建时为1,每次修改每次加1,作为一致性检查使用
	Contexts map[string]*Context // 上下文关系
	Data     []byte              // 数据，上层应用应该知道自己的到底是什么数据类型，并自行转换
}
