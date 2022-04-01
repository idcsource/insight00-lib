// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package dota

import (
	"fmt"

	"github.com/idcsource/insight00-lib/dotg/ns"
	"github.com/idcsource/insight00-lib/iendecode"
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

// 保持存活
func (c *Client) keepAlive() {
	c.iskeeplive = true
}

// 向服务器发送
// csend=Client Send，向服务器发送
// srb=Server Return Body，为了解开从服务器返回的数据
// 这个函数将会检查登录状态，如果服务器提示没有登录，将会尝试登录一次，然后再次发送要发送的数据
// 因此登录操作时这个不能用
func (c *Client) sendToServer(csend *ns.Client_Send, srb iendecode.BinaryCoder) (ssend *ns.Server_Send, err error) {
	ssend, err = c.doOneSend(csend, srb)
	if err != nil {
		return
	}
	if ssend.ReturnType == OPERATE_RETURN_ALL_OK {
		// 一切都OK
		return
	} else if ssend.ReturnType == OPERATE_RETURN_LOGIN_NO {
		// 没有登录，尝试重新登录并再次发送
		err = c.reLogin()
		if err != nil {
			return
		}
		ssend, err = c.doOneSend(csend, srb)
		return
	} else {
		return
	}
	return
}

// 与sendToServer配套的函数，打开一个连接，发送数据并回执，不做任何登录状态的检查
func (c *Client) doOneSend(csend *ns.Client_Send, srb iendecode.BinaryCoder) (ssend *ns.Server_Send, err error) {
	// 转换成byte
	csend_b, err := csend.ToBytes()
	if err != nil {
		return
	}
	// 打开连接
	cc, err := c.net_connect.OpenConnect()
	if err != nil {
		return
	}
	defer cc.Close()
	// 发送数据
	ssend_b, err := cc.SendAndReturn(csend_b)
	if err != nil {
		return
	}
	// 转码数据
	ssend = ns.New_Server_Send()
	err = ssend.FromBytes(ssend_b, srb)

	return
}

// 重新登录，与sendToServer配套的函数
func (c *Client) reLogin() (err error) {
	return
}

// 去进行登录
func (c *Client) ToLogin() (err error) {
	// 构建要发送的数据
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
	// 打开连接
	cc, err := c.net_connect.OpenConnect()
	defer cc.Close()
	// 发送并接收回执
	sre_b, err := cc.SendAndReturn(to_send_b)
	if err != nil {
		err = fmt.Errorf("dota: %v", err)
		return
	}
	// 解开回执
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

	// 启动保持存活
	go c.keepAlive()

	return
}
