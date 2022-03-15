// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package dot

import (
	"bytes"

	"github.com/idcsource/insight00-lib/iendecode"
)

// 上下文在dot中为map[string]*Context
// 每个dot支持多组上下文关系
// 每组上下文支持一个上文来源，和多个下文来源
type Context struct {
	Up   string            // 上下文的上游ID
	Down map[string]string // 上下文的下游ID
}

func (c *Context) MarshalBinary() (data []byte, err error) {
	b_buf := bytes.Buffer{}

	// 把Up转化为byte
	up_b := []byte(c.Up)
	up_b_len_b := iendecode.Uint64ToBytes(uint64(len(up_b)))
	b_buf.Write(up_b_len_b)
	b_buf.Write(up_b)

	// 把Down转化为byte
	down_b, err := iendecode.MapToBytes("map[string]string", c.Down)
	if err != nil {
		return
	}
	down_b_len_b := iendecode.Uint64ToBytes(uint64(len(down_b)))
	b_buf.Write(down_b_len_b)
	b_buf.Write(down_b)

	data = b_buf.Bytes()

	return
}

func (c *Context) UnmarshalBinary(data []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	c.Down = make(map[string]string)

	// 把up拿出来
	up_len := iendecode.BytesToUint64(data[0:8])
	c.Up = string(data[8 : 8+up_len])

	now_len := 8 + up_len

	// 把down拿出来
	down_len := iendecode.BytesToUint64(data[now_len : now_len+8])
	now_len = now_len + 8
	downi, err := iendecode.BytesToMap("map[string]string", data[now_len:now_len+down_len])
	if err != nil {
		return
	}
	c.Down = downi.(map[string]string)

	return
}

// 将上下文集合map[string]*Context转为byte
func ContextMapToByte(m map[string]*Context) (data []byte, err error) {
	b_buf := bytes.Buffer{}

	for key, _ := range m {
		key_b := []byte(key)
		key_len := iendecode.Uint64ToBytes(uint64(len(key_b)))
		b_buf.Write(key_len)
		b_buf.Write(key_b)

		c_b, err := m[key].MarshalBinary()
		if err != nil {
			break
		}
		c_b_len_b := iendecode.Uint64ToBytes(uint64(len(c_b)))
		b_buf.Write(c_b_len_b)
		b_buf.Write(c_b)
	}

	if err != nil {
		return
	}

	data = b_buf.Bytes()

	return
}

// 将byte转为上下文集合map[string]*Context
func ByteToContextMap(data []byte) (m map[string]*Context, err error) {
	m = make(map[string]*Context)

	var i uint64 = 0
	b_len := uint64(len(data))
	for {
		if i >= b_len {
			break
		}
		// key的长
		key_len_b := data[i:8]
		key_len := iendecode.BytesToUint64(key_len_b)
		i = i + 8

		// key
		key_b := data[i : i+key_len]
		key := string(key_b)
		i = i + key_len

		m[key] = &Context{}

		// 上下文的长
		c_b_len_b := data[i : i+8]
		c_b_len := iendecode.BytesToUint64(c_b_len_b)
		i = i + 8

		// 上下文
		c_b := data[i : i+c_b_len]
		i = i + c_b_len

		err = m[key].UnmarshalBinary(c_b)
		if err != nil {
			break
		}
	}

	if err != nil {
		return
	}

	return
}
