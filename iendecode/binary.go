// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package iendecode

import (
	"bytes"
	"encoding/binary"
)

func ToBinary(e interface{}) (b []byte, err error) {
	var bin_buf bytes.Buffer
	err = binary.Write(&bin_buf, binary.BigEndian, e)
	if err != nil {
		return
	}
	b = bin_buf.Bytes()
	return
}

func FromBinary(b []byte, e interface{}) (err error) {
	bin_buf := bytes.NewBuffer(b)
	err = binary.Read(bin_buf, binary.BigEndian, e)
	return
}
