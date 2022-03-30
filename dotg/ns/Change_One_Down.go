// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package ns

import (
	"bytes"

	"github.com/idcsource/insight00-lib/iendecode"
)

// 修改、删除、读取某个down都用这个
type Change_One_Down struct {
	DotName     string
	ContextName string
	DownName    string
	DownValue   string
}

func (c *Change_One_Down) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer

	// DotName
	dotn_b := []byte(c.DotName)
	dotn_b_l := len(dotn_b)
	dotn_b_l_b := iendecode.Uint64ToBytes(uint64(dotn_b_l))
	buf.Write(dotn_b_l_b)
	buf.Write(dotn_b)
	// ContextName
	cn_b := []byte(c.ContextName)
	cn_b_l := len(cn_b)
	cn_b_l_b := iendecode.Uint64ToBytes(uint64(cn_b_l))
	buf.Write(cn_b_l_b)
	buf.Write(cn_b)
	// DownName
	dn_b := []byte(c.DownName)
	dn_b_l := len(dn_b)
	dn_b_l_b := iendecode.Uint64ToBytes(uint64(dn_b_l))
	buf.Write(dn_b_l_b)
	buf.Write(dn_b)
	// DownValue
	dv_b := []byte(c.DownValue)
	dv_b_l := len(dv_b)
	dv_b_l_b := iendecode.Uint64ToBytes(uint64(dv_b_l))
	buf.Write(dv_b_l_b)
	if dv_b_l != 0 {
		buf.Write(dv_b)
	}

	data = buf.Bytes()

	return
}

func (c *Change_One_Down) UnmarshalBinary(data []byte) (err error) {
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
	c.DotName = string(name_b)
	// ContextName
	cn_b_len_b := buf.Next(8)
	cn_b_len := iendecode.BytesToUint64(cn_b_len_b)
	cn_b := buf.Next(int(cn_b_len))
	c.ContextName = string(cn_b)
	// DownName
	dn_b_len_b := buf.Next(8)
	dn_b_len := iendecode.BytesToUint64(dn_b_len_b)
	dn_b := buf.Next(int(dn_b_len))
	c.DownName = string(dn_b)
	// DownValue
	dv_b_len_b := buf.Next(8)
	dv_b_len := iendecode.BytesToUint64(dv_b_len_b)
	if dv_b_len != 0 {
		dv_b := buf.Next(int(dv_b_len))
		c.DownValue = string(dv_b)
	}

	return
}
