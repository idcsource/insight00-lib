// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package dota

import (
	"github.com/idcsource/insight00-lib/nst"
)

type Client struct {
	net_connect *nst.Client // 网络连接
	username    string      // 用户名
	password    string      // 密码
	unid        string      // 登陆后从服务器获得的
}
