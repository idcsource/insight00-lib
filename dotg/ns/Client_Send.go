// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package ns

import (
	"bytes"
	"time"

	"github.com/idcsource/insight00-lib/iendecode"
)

// 这是客户端往服务器发送的完整信息
type Client_Send struct {
	Time         time.Time        // 触发操作的时间
	OperateType  uint8            // 操作类型,OPERATE_TYPE_开头的那些
	LoginInfo    *Login_Base_Info // 登陆信息，这个有可能是空的
	OperateBlock string           // 要操作的block，这个可能为空
	OperateBody  []byte           // 操作的具体请求体
}

func New_Client_Send() (c *Client_Send) {
	c = &Client_Send{
		Time:         time.Now(),
		OperateType:  0,
		LoginInfo:    &Login_Base_Info{},
		OperateBlock: "",
	}
	return
}

func (c *Client_Send) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer

	// Time（15byte）
	time_b, err := c.Time.MarshalBinary()
	if err != nil {
		return
	}
	buf.Write(time_b)
	// OperateType(1byte)
	opt_b := iendecode.Uint8ToBytes(c.OperateType)
	buf.Write(opt_b)
	// LoginInfo
	li_b, err := c.LoginInfo.MarshalBinary()
	if err != nil {
		return
	}
	li_b_l := len(li_b)
	li_b_l_b := iendecode.Uint64ToBytes(uint64(li_b_l))
	buf.Write(li_b_l_b)
	if li_b_l != 0 {
		buf.Write(li_b)
	}
	//OperateBlock
	o_b := []byte(c.OperateBlock)
	o_b_len := len(o_b)
	o_b_len_b := iendecode.Uint64ToBytes(uint64(o_b_len))
	buf.Write(o_b_len_b)
	if o_b_len != 0 {
		buf.Write(o_b)
	}
	// OperateBody
	ob_b_l := len(c.OperateBody)
	ob_b_l_b := iendecode.Uint64ToBytes(uint64(ob_b_l))
	buf.Write(ob_b_l_b)
	if ob_b_l != 0 {
		buf.Write(c.OperateBody)
	}

	data = buf.Bytes()

	return
}

func (c *Client_Send) UnmarshalBinary(data []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	buf := bytes.NewBuffer(data)

	// Time（15byte）
	time_b := buf.Next(15)
	err = c.Time.UnmarshalBinary(time_b)
	if err != nil {
		return
	}
	// OperateType(1byte)
	ot_b := buf.Next(1)
	c.OperateType = iendecode.BytesToUint8(ot_b)
	// LoginInfo
	li_b_l_b := buf.Next(8)
	li_b_l := iendecode.BytesToUint64(li_b_l_b)
	if li_b_l != 0 {
		li_b := buf.Next(int(li_b_l))
		c.LoginInfo = &Login_Base_Info{}
		err = c.LoginInfo.UnmarshalBinary(li_b)
		if err != nil {
			return
		}
	}
	//OperateBlock
	o_b_l_b := buf.Next(8)
	o_b_l := iendecode.BytesToUint64(o_b_l_b)
	if o_b_l != 0 {
		o_b := buf.Next(int(o_b_l))
		c.OperateBlock = string(o_b)
	}
	// OperateBody
	ob_b_l_b := buf.Next(8)
	ob_b_l := iendecode.BytesToUint64(ob_b_l_b)
	if ob_b_l != 0 {
		c.OperateBody = buf.Next(int(ob_b_l))
	}

	return
}
