// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ]
// Stephen Fire Meditation Qin [ 火志溟 ] -> stephenfmqin@gmail.com
// This source code is governed by GNU LGPL v3 license

package dot

// 上下文在dot中为map[string]*Context
// 每个dot支持多组上下文关系
// 每组上下文支持一个上文来源，和多个下文来源
type Context struct {
	Version uint8             // 涉及所操作dot的版本号，现在已经为V2
	Up      string            // 上下文的上游ID
	UpData  []byte            // 上下文的上游配置数据
	Down    map[string][]byte // 上下文的下游ID，以及携带的配置数据
}

// 上下文关系索引，对应文件：(dot id的40位sha1散列)_context_index
type ContextIndex struct {
	Status      _DotContextIndex_Status // 标记的状态
	ContextName string                  // 上下文关系的名称
}

// Context中Down的状态
type ContextDownStatus struct {
	HardIndex uint64                        //索引的物理索引编号（不管是否删除）
	HardCount uint64                        // 索引的物理起始字节
	Name      string                        // Down的名字
	Status    _DotContextUpDownIndex_Status // 状态
	DataLen   uint64                        // 数据长度
}
