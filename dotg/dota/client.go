// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package dota

import (
	"fmt"

	"github.com/idcsource/insight00-lib/dotg/ns"
	"github.com/idcsource/insight00-lib/nst"
)

type Client struct {
	net_connect    *nst.Client // 网络连接
	username       string      // 用户名
	password       string      // 密码
	unid           string      // 登录后从服务器获得的
	islogin        bool        // 是否登录了
	iskeeplive     bool        // 是否保持存活
	keepalive_chan chan bool   // 保持存活的chan状态
}

func NewClient(username string, password string, net_connect *nst.Client) (client *Client) {
	client = &Client{
		username:       username,
		password:       password,
		net_connect:    net_connect,
		islogin:        false,
		iskeeplive:     false,
		keepalive_chan: make(chan bool),
	}
	return
}

// 返回错误检查
func (c *Client) checkErr(sre *ns.Server_Send) (err error) {
	err = fmt.Errorf("dota: %v", sre.ReturnErr)
	return
}

// 去进行登录
func (c *Client) ToLogin() (err error) {
	cc, err := c.net_connect.OpenConnect()
	defer cc.Close()

	to_login := &ns.To_Login{
		Name:     c.username,
		Password: c.password,
	}
	to_send := ns.New_Client_Send()
	to_send.OperateBody = to_login
	to_send.OperateType = OPERATE_TYPE_LOGIN
	to_send_b, err := to_send.ToBytes()
	if err != nil {
		err = fmt.Errorf("dota: %v", err)
		return
	}

	sre_b, err := cc.SendAndReturn(to_send_b)
	sre_body := &ns.String_In_Byte{}
	sre := ns.New_Server_Send()
	err = sre.FromBytes(sre_b, sre_body)
	if err != nil {
		err = fmt.Errorf("dota: %v", err)
		return
	}
	if sre.ReturnType == OPERATE_RETURN_ALL_OK {
		c.islogin = true
		c.unid = sre_body.String
	} else {
		err = c.checkErr(sre) // 这里的错误先简单处理一下吧
		return
	}

	return
}
