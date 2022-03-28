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
	UserName    string // 用户名
	OldPassword string // sha1之后的老密码
	NewPassword string // Sha1之后的新密码
}

func (c *Change_Password) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer

	user_b := []byte(c.UserName)
	user_b_len := len(user_b)
	user_b_len_b := iendecode.Uint64ToBytes(uint64(user_b_len))
	buf.Write(user_b_len_b)
	buf.Write(user_b)

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

	user_b_len_b := buf.Next(8)
	user_b_len := iendecode.BytesToUint64(user_b_len_b)
	user_b := buf.Next(int(user_b_len))
	c.UserName = string(user_b)

	oldpassword_b := buf.Next(40)
	c.OldPassword = string(oldpassword_b)

	newpassword_b := buf.Next(40)
	c.NewPassword = string(newpassword_b)

	return
}
