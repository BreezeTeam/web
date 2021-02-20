package web

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isFuzzy  bool    // 是否模糊匹配，part 含有 : 或 * 时为true
}

//toString
func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isFuzzy)
}
//使用bfs来寻找parts对应的node
func (n *node)search(parts []string,height int) *node{
	//递归结束条件,这里不能直接取第一个字符，因为有可能第一个字符不存在
	//如果是*开始的，那么就可以把这里设置为一个节点，之后只要匹配到这节点的，都能匹配成功
	if len(parts) == height || strings.HasPrefix(n.part, "*"){
		//如果pattern，说明，这个路径没对应的method，匹配失败
		if n.pattern == ""{
			return nil
		}
		return n
	}
	part:=parts[height]
	children:=n.matchChildren(part)
	//循环该节点的每一个子节点，其中，对每一个节点再进行搜索
	for _, child := range children {
		result :=child.search(parts,height+1)
		if result!=nil{
			return result
		}
	}
	return nil
}



//插入
func (n *node) insert(pattern string, parts []string, height int) {
	//递归结束，插入完成
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	//如果children为空，表示该节点还没有对应的子节点，需要进行创建
	if child == nil {
		child = &node{part: part, isFuzzy: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	//进行下一层的插入
	child.insert(pattern,parts,height+1)
}

//通过part匹配节点的第一个子节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isFuzzy {
			return child
		}
	}
	return nil
}

//通过part匹配到全部的子节点
func (n *node) matchChildren(part string) []*node {
	children := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isFuzzy {
			children = append(children, child)
		}
	}
	return children
}
