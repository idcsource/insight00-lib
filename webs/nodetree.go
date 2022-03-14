// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package webs

import (
	"github.com/idcsource/insight00-lib/jconf"
)

// 创建根节点，内部调用使用，不需要直接新建
func AddRootNode(f FloorInterface, config *jconf.JsonConf) (root *NodeTree) {
	root = &NodeTree{
		name:        "Index",
		mark:        "",
		config:      config,
		if_children: false,
		node_type:   NODE_IS_ROOT,
		children:    make(map[string]*NodeTree),
		floor:       f,
	}
	return
}

// 增加一个普通的节点
func (nt *NodeTree) AddNode(name, mark string, f FloorInterface, config *jconf.JsonConf) (child *NodeTree) {
	nt.if_children = true
	child = &NodeTree{
		name:        name,
		mark:        mark,
		config:      config,
		if_children: false,
		node_type:   NODE_IS_NOMAL,
		children:    make(map[string]*NodeTree),
		floor:       f,
	}
	nt.children[mark] = child
	return
}

// 增加一个节点Door
func (nt *NodeTree) AddDoor(name, mark string, f FloorDoorInterface, config *jconf.JsonConf) (child *NodeTree) {
	nt.if_children = true
	child = &NodeTree{
		name:        name,
		mark:        mark,
		config:      config,
		if_children: false,
		node_type:   NODE_IS_DOOR,
		door:        f,
		children:    make(map[string]*NodeTree),
	}
	nt.children[mark] = child
	return
}

// 增加一个处理静态文件的节点（静态文件节点只能做最后一个节点，所以不再提供节点树的返回）
func (nt *NodeTree) AddStatic(mark, path string) {
	nt.if_children = true
	child := &NodeTree{
		name:        mark,
		mark:        mark,
		if_children: false,
		node_type:   NODE_IS_STATIC,
		floor: &StaticFileFloor{
			path: path,
		},
	}
	nt.children[mark] = child
}

// 增加一个空节点
func (nt *NodeTree) AddEmpty(name, mark string) (child *NodeTree) {
	nt.if_children = true
	child = &NodeTree{
		name:        name,
		mark:        mark,
		if_children: false,
		node_type:   NODE_IS_EMPTY,
		floor:       &EmptyFloor{},
		children:    make(map[string]*NodeTree),
	}
	nt.children[mark] = child
	return
}

// （此方法不需要手动调用）找到现在需要去运行的Floor，如果找到了则返回的NodeTree为找到的NodeTree，如果找不到，则返回的是能找到的最后一个的NodeTree
func (nt *NodeTree) getRunFloor(rt Runtime) (nnt *NodeTree, rtn Runtime, nothing bool) {
	nothing = false
	var ntok bool
	nnt, ntok = nt.children[rt.NowRoutePath[0]]
	if ntok == false {
		// 如果找不到，就启动nothing，最终也就是走404
		nothing = true
		rtn = rt
		nnt = nt
		//rtn.ConfigTree[nnt.Mark] = nnt.Config;
		rtn.MyConfig = nnt.config
	} else {
		rt.NowRoutePath = rt.NowRoutePath[1:]      //先将总的请求减去一个
		rt.RealNode = rt.RealNode + "/" + nnt.mark //将RealNode修改掉
		nothing = false
		if nnt.if_children == false || len(rt.NowRoutePath) == 0 {
			//如果新的节点没有子节点存在，或者NowRoutePath已经为空，则不再继续下面的操作
			rtn = rt
			//rtn.ConfigTree[nnt.Mark] = nnt.Config;
			rtn.MyConfig = nnt.config
			return
		} else {
			nnt, rtn, nothing = nnt.getRunFloor(rt)
			return
		}
	}
	return
}
