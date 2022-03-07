// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package logs

// 新建一个运行时日志
func NewRuntimeLog(maxnum int) (rtl *RuntimeLog) {
	var logs []string
	if maxnum > 0 {
		logs = make([]string, 0, maxnum)
	}
	rtl = &RuntimeLog{
		logs:   logs,
		maxnum: maxnum,
	}
	return
}

// 设置运行时日志的最大条目数
func (rtl *RuntimeLog) SetMax(maxnum int) (err error) {
	rtl.maxnum = maxnum
	return
}

// 写入一条运行时日志
func (rtl *RuntimeLog) Write(p []byte) (n int, err error) {
	if rtl.maxnum > 0 {
		llen := len(rtl.logs)
		if llen >= rtl.maxnum {
			logs := make([]string, 0, rtl.maxnum)
			logs = rtl.logs[1:]
			logs = append(logs, string(p))
			rtl.logs = logs
		} else {
			rtl.logs = append(rtl.logs, string(p))
		}
	} else {
		rtl.logs = append(rtl.logs, string(p))
	}
	n = len(string(p))
	return
}

// 输出运行时日志
func (rtl *RuntimeLog) Output() (output []string) {
	return rtl.logs
}
