// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package iendecode

// BinaryCoder 其他数据类型与[]byte间的转换接口
type BinaryCoder interface {
	MarshalBinary() (data []byte, err error)
	UnmarshalBinary(data []byte) error
}
