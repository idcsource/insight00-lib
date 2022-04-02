// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package ns

import (
	"bytes"

	"github.com/idcsource/insight00-lib/iendecode"
)

// 这是客户端往服务器发送的完整信息
type Server_Send struct {
	ReturnType uint8  // 操作类型，OPERATE_RETURN_的那些
	ReturnErr  string // 返回的错误，这个可能为空
	ReturnBody []byte // 返回的结构体，可能为空
}

func New_Server_Send() (c *Server_Send) {
	c = &Server_Send{
		ReturnType: 0,
		ReturnErr:  "",
	}
	return
}

func (c *Server_Send) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer

	// ReturnType(1byte)
	opt_b := iendecode.Uint8ToBytes(c.ReturnType)
	buf.Write(opt_b)
	//ReturnErr
	o_b := []byte(c.ReturnErr)
	o_b_len := len(o_b)
	o_b_len_b := iendecode.Uint64ToBytes(uint64(o_b_len))
	buf.Write(o_b_len_b)
	if o_b_len != 0 {
		buf.Write(o_b)
	}
	// ReturnBody
	ob_b_l := len(c.ReturnBody)
	ob_b_l_b := iendecode.Uint64ToBytes(uint64(ob_b_l))
	buf.Write(ob_b_l_b)
	if ob_b_l != 0 {
		buf.Write(c.ReturnBody)
	}

	data = buf.Bytes()

	return
}

func (c *Server_Send) UnmarshalBinary(data []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	buf := bytes.NewBuffer(data)

	// ReturnType(1byte)
	ot_b := buf.Next(1)
	c.ReturnType = iendecode.BytesToUint8(ot_b)
	//ReturnErr
	o_b_l_b := buf.Next(8)
	o_b_l := iendecode.BytesToUint64(o_b_l_b)
	if o_b_l != 0 {
		o_b := buf.Next(int(o_b_l))
		c.ReturnErr = string(o_b)
	}
	// ReturnBody
	ob_b_l_b := buf.Next(8)
	ob_b_l := iendecode.BytesToUint64(ob_b_l_b)
	if ob_b_l != 0 {
		c.ReturnBody = buf.Next(int(ob_b_l))
	}

	return
}
