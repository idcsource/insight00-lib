// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package ns

import (
	"bytes"

	"github.com/idcsource/insight00-lib/iendecode"
)

// 这是dota客户端进行登陆的数据结构
type Change_Password struct {
	LoginInfo   *Login_Base_Info
	OldPassword string // sha1之后的老密码
	NewPassword string // Sha1之后的新密码
}

func (c *Change_Password) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer

	login_info_b, err := c.LoginInfo.MarshalBinary()
	if err != nil {
		return
	}
	login_info_b_len := len(login_info_b)
	login_info_b_len_b := iendecode.Uint64ToBytes(uint64(login_info_b_len)) // 8位的长度

	buf.Write(login_info_b_len_b)
	buf.Write(login_info_b)

	oldpassword_b := []byte(c.OldPassword)
	buf.Write(oldpassword_b)

	newpassword_b := []byte(c.NewPassword)
	buf.Write(newpassword_b)

	data = buf.Bytes()
	return
}

func (c *Change_Password) UnmarshalBinary(data []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	buf := bytes.NewBuffer(data)

	login_info_b_len_b := buf.Next(8)
	login_info_b_len := iendecode.BytesToUint64(login_info_b_len_b)
	login_info_b := buf.Next(int(login_info_b_len))
	c.LoginInfo = &Login_Base_Info{}
	err = c.LoginInfo.UnmarshalBinary(login_info_b)
	if err != nil {
		return
	}

	oldpassword_b := buf.Next(40)
	c.OldPassword = string(oldpassword_b)

	newpassword_b := buf.Next(40)
	c.NewPassword = string(newpassword_b)

	return
}
