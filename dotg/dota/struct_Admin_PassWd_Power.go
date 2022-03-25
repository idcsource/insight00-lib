// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package dota

import (
	"fmt"

	"github.com/idcsource/insight00-lib/iendecode"
)

// 用户的密码和权限类型数据
type Admin_PassWd_Power struct {
	Password  string // 密码的sha1
	PowerType uint8  // 用户权限
}

func (app *Admin_PassWd_Power) MarshalBinary() (data []byte, err error) {
	pwd_sha1_b := []byte(app.Password)
	power_b := iendecode.Uint8ToBytes(app.PowerType)
	data = append(pwd_sha1_b, power_b...)
	return
}

func (app *Admin_PassWd_Power) UnmarshalBinary(data []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	if len(data) != 41 {
		err = fmt.Errorf("This is not a Admin_PassWd_Power")
		return
	}
	pwd_sha1_b := data[0:40]
	app.Password = string(pwd_sha1_b)
	power := iendecode.BytesToUint8(data[40:41])
	app.PowerType = power
	return
}
