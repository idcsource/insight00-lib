// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package dota

import (
	"fmt"
	"sync"
	"time"

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
	loged_user      map[string]logedUser   // 已经登陆用户，string为用户名
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
	if cs.OperateType == OPERATE_TYPE_LOGIN {
		// 这里是执行登录
	} else {
		// 这里加登录判断，如果没有登录就不用执行下面的工作了
		switch cs.OperateType {
		case OPERATE_TYPE_KEEPLIVE:
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
		}
	}
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
