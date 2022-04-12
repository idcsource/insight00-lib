// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package dota

import (
	"fmt"
	"sync"
	"time"

	"github.com/idcsource/insight00-lib/base"
	"github.com/idcsource/insight00-lib/dotg/dot"
	"github.com/idcsource/insight00-lib/dotg/ns"
	"github.com/idcsource/insight00-lib/logs"
	"github.com/idcsource/insight00-lib/nst"
)

type Server struct {
	user_lock       *sync.RWMutex          // 操作user时的锁
	block_lock      *sync.RWMutex          // 操作block时的锁
	run_log         *logs.Logs             // 运行日志
	err_log         *logs.Logs             // 错误日志
	dota_op         *dot.DotsOp            // Dot-Area默认block的操作
	loged_user      map[string]*logedUser  // 已经登陆用户，string为用户名
	block_op        map[string]*dot.DotsOp // 所有管理block的操作，将一个op加入时，需要启动block_lock的锁
	run_count       uint64                 // 操作计数
	closed          bool                   // 关闭状态
	closeing_signal chan bool              // 关闭信号
	run_wait        *sync.WaitGroup        // 执行计数
}

// 已经登陆的用户
type logedUser struct {
	username  string
	unid      map[string]time.Time // string为unid，time则为活动时间，目的是一个用户可以多次登陆
	authority uint8                // 用户权限，管理员还是一般
	wrable    map[string]bool      // string是能访问的block的名字
	lock      *sync.RWMutex        // 锁，当修改unid的时候
}

// nst的ConnExecer实现
func (s *Server) NSTexec(ce *nst.ConnExec) (stat nst.SendStat, err error) {
	if s.closed == true {
		// 这里发送服务器关闭的状态
		return
	}
	s.run_wait.Add(1)
	defer s.run_wait.Done()

	// 接收数据
	cs_b, err := ce.GetData()
	if err != nil {
		err = fmt.Errorf("dota: %v", err)
		return
	}
	// 解码数据
	cs := ns.New_Client_Send()
	err = cs.UnmarshalBinary(cs_b)
	if err != nil {
		err = fmt.Errorf("dota: %v", err)
		return
	}
	var ss *ns.Server_Send
	if cs.OperateType == OPERATE_TYPE_LOGIN {
		// 这里是执行登录
		ss, err = s.doLogin(cs)
	} else {
		// 这里加登录判断，如果没有登录就不用执行下面的工作了
		switch cs.OperateType {
		case OPERATE_TYPE_KEEPLIVE:
			ss, err = s.doKeepLive(cs)
		case OPERATE_TYPE_CHANGE_PASSWORD:
		case OPERATE_TYPE_NEW_USER:
		case OPERATE_TYPE_USER_ADD_BLOCK:
		case OPERATE_TYPE_USER_DEL_BLOCK:
		case OPERATE_TYPE_DEL_USER:
		case OPERATE_TYPE_NEW_BLOCK:
		case OPERATE_TYPE_DEL_BLOCK:
		case OPERATE_TYPE_NEW_DOT:
		case OPERATE_TYPE_NEW_DOT_WITH_CONTEXT:
		case OPERATE_TYPE_DEL_DOT:
		case OPERATE_TYPE_UPDATE_DATA:
		case OPERATE_TYPE_READ_DATA:
		case OPERATE_TYPE_UPDATE_ONE_DOWN:
		case OPERATE_TYPE_UPDATE_ONE_UP:
		case OPERATE_TYPE_DEL_ONE_DOWN:
		case OPERATE_TYPE_ADD_CONTEXT:
		case OPERATE_TYPE_UPDATE_CONTEXT:
		case OPERATE_TYPE_DEL_CONTEXT:
		case OPERATE_TYPE_READ_CONTEXT:
		case OPERATE_TYPE_READ_ONE_UP:
		case OPERATE_TYPE_READ_ONE_DOWN:
		case OPERATE_TYPE_READ_DATA_TV:
		case OPERATE_TYPE_READ_INDEX_TV:
		case OPERATE_TYPE_READ_CONTEXT_TV:
		default:
			// 这里是没有任何知道的请求
			ss, err = s.doNoOperateType()
		}
	}
	if err != nil {
		// 这里负责所有没有被提前返回的错误
		ss, err = s.doAllErr(err)
	}
	// 如果没错误，就负责发送Server_Send
	ss_b, _ := ss.MarshalBinary()
	err = ce.SendData(ss_b)

	return
}

// 当执行关闭方法时Close()的时候，首先将closed置于true，这样NSTexec将不在接收新的请求
// 执行WaitGroup的Wait进行等待堵塞，预防还有未结束的处理
// 当WaitGroup不再堵塞，则执行完Close
// 建议先执行nst.Server的Close，让其不在接收新请求，再执行这个Close将所有剩下的执行完毕
func (s *Server) Close() (err error) {
	s.closed = true
	s.run_wait.Wait()
	return

}

// 处理所有可能的错误
func (s *Server) doAllErr(err error) (ss *ns.Server_Send, errs error) {
	// 构建个错误，发出去
	ss = &ns.Server_Send{
		ReturnType: OPERATE_RETURN_ERROR,
		ReturnErr:  fmt.Sprint(err),
	}

	return
}

// 执行如果没有对应操作的请求
func (s *Server) doNoOperateType() (ss *ns.Server_Send, err error) {
	// 构建个错误，发出去
	ss = &ns.Server_Send{
		ReturnType: OPERATE_RETURN_TYPE_NOT_HAVE,
		ReturnErr:  "The Opereate Not Exist",
	}

	return
}

// 执行登录
func (s *Server) doLogin(cs *ns.Client_Send) (ss *ns.Server_Send, err error) {
	// 解开To_Login
	csb := &ns.To_Login{}
	err = csb.UnmarshalBinary(cs.OperateBody)
	if err != nil {
		return
	}
	// 加操作User的所
	s.user_lock.Lock()
	defer s.user_lock.Unlock()

	ss = &ns.Server_Send{}

	// 查看是不是有这个dot
	have, err := s.dota_op.HaveDot(DEFAULT_USER_PREFIX + csb.Name)
	if err != nil {
		return
	}
	if have == false {
		// 如果没有的话，发送用户名和密码错误
		ss.ReturnType = OPERATE_RETURN_PASSWD_NO
	} else {
		// 去默认的block中找这个user
		var user_b []byte
		user_b, _, err = s.dota_op.ReadData(DEFAULT_USER_PREFIX + csb.Name)
		if err != nil {
			return
		}
		// 解开用户信息
		users := &ns.User_PassWd_Power{}
		err = users.UnmarshalBinary(user_b)
		if err != nil {
			return
		}
		// 对比密码
		if csb.Password == users.Password {
			// 先找以下是否已经登过
			if _, exit := s.loged_user[csb.Name]; exit == false {
				// 如果没有登过，就构建一个
				s.loged_user[csb.Name] = &logedUser{
					username:  csb.Name,
					unid:      make(map[string]time.Time),
					authority: USER_AUTHORITY_NO,
					wrable:    make(map[string]bool),
					lock:      new(sync.RWMutex),
				}
			}
			// 查看权限
			if users.PowerType == USER_AUTHORITY_ADMIN {
				// 如果是管理员
				s.loged_user[csb.Name].authority = USER_AUTHORITY_ADMIN
			} else {
				// 如果不是管理员，找block
				var blocks *dot.Context
				blocks, err = s.dota_op.ReadContext(csb.Name, "block")
				if err != nil {
					return
				}
				s.loged_user[csb.Name].wrable = make(map[string]bool)
				for k, _ := range blocks.Down {
					s.loged_user[csb.Name].wrable[k] = true
				}
				s.loged_user[csb.Name].authority = USER_AUTHORITY_NORMAL
			}
			// 清理过期unid
			s.cleanUnid(csb.Name)
			// 生成唯一unid，并写进去
			unid := base.Unid(1, csb.Name)
			s.loged_user[csb.Name].unid[unid] = time.Now()
			// 加入loged_user
		} else {
			// 如果密码不对，发送用户名和密码错误
			ss.ReturnType = OPERATE_RETURN_PASSWD_NO
		}
	}
	return
}

// 清理过期的unid
func (s *Server) cleanUnid(name string) {
	for k, v := range s.loged_user[name].unid {
		if v.Unix()+SERVER_OUTLOG_TIME < time.Now().Unix() {
			delete(s.loged_user[name].unid, k)
		}
	}
}

// 执行续期
func (s *Server) doKeepLive(cs *ns.Client_Send) (ss *ns.Server_Send, err error) {
	// 解开数据
	loginbaseinfo := &ns.Login_Base_Info{}
	err = loginbaseinfo.UnmarshalBinary(cs.OperateBody)
	if err != nil {
		return
	}
	// 对用户项目加锁
	s.user_lock.Lock()
	defer s.user_lock.Unlock()

	ss = &ns.Server_Send{}

	// 找找有没有这个用户信息
	if _, have := s.loged_user[loginbaseinfo.Name]; have == false {
		ss.ReturnType = OPERATE_RETURN_LOGIN_NO
		return
	}
	if _, have := s.loged_user[loginbaseinfo.Name].unid[loginbaseinfo.Unid]; have == false {
		ss.ReturnType = OPERATE_RETURN_LOGIN_NO
		return
	}
	// 看有无过期
	if s.loged_user[loginbaseinfo.Name].unid[loginbaseinfo.Unid].Unix()+SERVER_OUTLOG_TIME < time.Now().Unix() {
		delete(s.loged_user[loginbaseinfo.Name].unid, loginbaseinfo.Unid)
		ss.ReturnType = OPERATE_RETURN_LOGIN_NO // 也是没登录
		return
	}

	// 都正常
	s.loged_user[loginbaseinfo.Name].unid[loginbaseinfo.Unid] = time.Now()
	ss.ReturnType = OPERATE_RETURN_ALL_OK

	return
}
