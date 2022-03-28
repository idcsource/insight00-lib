// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package ns

import (
	"bytes"

	"github.com/idcsource/insight00-lib/iendecode"
)

// 给用户添加或删除某个block权限
type User_Block struct {
	UserName  string
	BlockName string
}

func (u *User_Block) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer

	user_b := []byte(u.UserName)
	user_b_len := len(user_b)
	user_b_len_b := iendecode.Uint64ToBytes(uint64(user_b_len)) // 8位的长度
	buf.Write(user_b_len_b)
	buf.Write(user_b)

	block_b := []byte(u.BlockName)
	block_b_len := len(block_b)
	block_b_len_b := iendecode.Uint64ToBytes(uint64(block_b_len))
	buf.Write(block_b_len_b)
	buf.Write(block_b)

	data = buf.Bytes()
	return
}

func (u *User_Block) UnmarshalBinary(data []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	buf := bytes.NewBuffer(data)

	user_b_len_b := buf.Next(8)
	user_b_len := iendecode.BytesToUint64(user_b_len_b)
	user_b := buf.Next(int(user_b_len))
	u.UserName = string(user_b)

	block_b_len_b := buf.Next(8)
	block_b_len := iendecode.BytesToUint64(block_b_len_b)
	block_b := buf.Next(int(block_b_len))
	u.BlockName = string(block_b)

	return
}
