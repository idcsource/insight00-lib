// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package ns

// 主要是为了实现iendecode的接口
type String_In_Byte struct {
	String string
}

func (s *String_In_Byte) MarshalBinary() (data []byte, err error) {

	data = []byte(s.String)

	return
}

func (s *String_In_Byte) UnmarshalBinary(data []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	s.String = string(data)

	return
}
