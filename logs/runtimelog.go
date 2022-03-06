// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package logs

func NewRuntimeLogContent(maxnum int) (rtl *RuntimeLogContent) {
	var logs []string
	if maxnum > 0 {
		logs = make([]string, 0, maxnum)
	}
	rtl = &RuntimeLogContent{
		logs:   logs,
		maxnum: maxnum,
	}
	return
}

func (rtl *RuntimeLogContent) Write(p []byte) (n int, err error) {
	if rtl.maxnum > 0 {
		llen := len(rtl.logs)
		if llen >= rtl.maxnum {
			logs := make([]string, 0, maxnum)
			logs = rtl.logs[1:]
			logs = append(logs, string(p))
			rtl.logs = logs
		} else {
			rtl.logs = append(rtl.logs, string(p))
		}
	} else {
		rtl.logs = append(rtl.logs, string(p))
	}
	return
}

func (rtl *RuntimeLogContent) Output() (output []string) {
	return rtl.logs
}
