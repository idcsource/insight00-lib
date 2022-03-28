// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package ns

import (
	"bytes"

	"github.com/idcsource/insight00-lib/iendecode"
)

// 用户的密码和权限类型数据，这个是在server内使用的
type User_PassWd_Power struct {
	Name      string // 用户名
	Password  string // 密码的sha1
	PowerType uint8  // 用户权限
}

func (u *User_PassWd_Power) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer

	name_b := []byte(u.Name)
	name_b_len := len(name_b)
	name_b_len_b := iendecode.Uint64ToBytes(uint64(name_b_len))
	buf.Write(name_b_len_b)
	buf.Write(name_b)

	password_b := []byte(u.Password)
	buf.Write(password_b)

	pt_b := iendecode.Uint8ToBytes(u.PowerType)
	buf.Write(pt_b)

	data = buf.Bytes()

	return
}

func (u *User_PassWd_Power) UnmarshalBinary(data []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	buf := bytes.NewBuffer(data)

	name_b_len_b := buf.Next(8)
	name_b_len := iendecode.BytesToUint64(name_b_len_b)
	name_b := buf.Next(int(name_b_len))
	u.Name = string(name_b)

	password_b := buf.Next(40)
	u.Password = string(password_b)

	pt_b := buf.Next(1)
	u.PowerType = iendecode.BytesToUint8(pt_b)

	return
}
