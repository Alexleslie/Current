package gee

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string  //待匹配路由，例如/p/:lang，是否为一个完整URL，不是则为"”
	part     string  //路由中的一部分，例如：lang，URL切割后的块值
	children []*node //该节点的子节点，例如[doc,tutorial,intro]
	isWild   bool    //是否模糊匹配，part含有；或*时候为true
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// 找第一个匹配成功的节点，用于插入（找到一个就立即返回），如果为模糊匹配也可以成功匹配上
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 匹配所有成功的节点后返回，用于查找，必须返回所有可能成功的子节点来进行遍历查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 先匹配后插入
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 查找匹配的URL
func (n *node) search(parts []string, height int) *node {
	//递归终止条件：找到末尾了或者通配符
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			//pattern为空字符串则表示它不是一个完整的URL，匹配失败
			return nil
		}
		return n
	}

	part := parts[height]
	//遍历查找所有可能的子路径
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
