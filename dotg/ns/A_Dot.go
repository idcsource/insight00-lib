// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package ns

import (
	"bytes"

	"github.com/idcsource/insight00-lib/iendecode"
)

// 这是添加一个dot时使用
type A_Dot struct {
	DotName string
	DotData []byte
}

func (a *A_Dot) MarshalBinary(data []byte, err error) {
	var buf bytes.Buffer

	// DotName
	name_b := []byte(a.DotName)
	name_b_len := len(name_b)
	name_b_len_b := iendecode.Uint64ToBytes(uint64(name_b_len))
	buf.Write(name_b_len_b)
	buf.Write(name_b)
	// DotData
	data_len := len(a.DotData)
	data_len_b := iendecode.Uint64ToBytes(uint64(data_len))
	buf.Write(data_len_b)
	if data_len != 0 {
		buf.Write(a.DotData)
	}

	data = buf.Bytes()

	return
}

func (a *A_Dot) UnmarshalBinary(data []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	buf := bytes.NewBuffer(data)
	// DotName
	name_b_len_b := buf.Next(8)
	name_b_len := iendecode.BytesToUint64(name_b_len_b)
	name_b := buf.Next(int(name_b_len))
	a.DotName = string(name_b)
	// DotData
	data_len_b := buf.Next(8)
	data_len := iendecode.BytesToUint64(data_len_b)
	if data_len != 0 {
		a.DotData = buf.Next(int(data_len))
	}

	return
}
