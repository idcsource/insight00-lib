// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package ns

import (
	"bytes"

	"github.com/idcsource/insight00-lib/dotg/dot"
	"github.com/idcsource/insight00-lib/iendecode"
)

// 这是添加一个dot时使用，如果不添加context，则可以忽略context的数据
type A_Dot_With_A_Context struct {
	DotName     string
	DotData     []byte
	ContextName string
	Context     *dot.Context
}

func (a *A_Dot_With_A_Context) MarshalBinary(data []byte, err error) {
	return
}

func (a *A_Dot_With_A_Context) UnmarshalBinary(data []byte) (err error) {
	return
}
