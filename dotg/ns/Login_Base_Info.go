// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package ns

import (
	"bytes"

	"github.com/idcsource/insight00-lib/iendecode"
)

// 这是dota客户端登录后，每次操作客户端都要发给服务器的
type Login_Base_Info struct {
	Name string
	Unid string
}

func (l *Login_Base_Info) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer

	name_b := []byte(l.Name)
	name_b_len := len(name_b)
	name_b_len_b := iendecode.Uint64ToBytes(uint64(name_b_len)) // 8位的长度
	uuid_b := []byte(l.Unid)

	buf.Write(name_b_len_b)
	buf.Write(name_b)
	buf.Write(uuid_b)

	data = buf.Bytes()
	return
}

func (l *Login_Base_Info) UnmarshalBinary(data []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	buf := bytes.NewBuffer(data)
	name_b_len_b := buf.Next(8)
	name_b_len := iendecode.BytesToUint64(name_b_len_b)
	name_b := buf.Next(int(name_b_len))
	l.Name = string(name_b)

	uuid_b := buf.Next(40)
	l.Unid = string(uuid_b)

	return
}
