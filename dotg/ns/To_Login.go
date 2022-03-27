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
type To_Login struct {
	Name     string
	Password string // sha1之后的
}

func (t *To_Login) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer

	name_b := []byte(t.Name)
	name_b_len := len(name_b)
	name_b_len_b := iendecode.Uint64ToBytes(uint64(name_b_len)) // 8位的长度
	passwd_b := []byte(t.Password)

	buf.Write(name_b_len_b)
	buf.Write(name_b)
	buf.Write(passwd_b)

	data = buf.Bytes()
	return
}

func (t *To_Login) UnmarshalBinary(data []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	buf := bytes.NewBuffer(data)
	name_b_len_b := buf.Next(8)
	name_b_len := iendecode.BytesToUint64(name_b_len_b)
	name_b := buf.Next(int(name_b_len))
	t.Name = string(name_b)

	passwd_b := buf.Next(40)
	t.Password = string(passwd_b)

	return
}
