package gee

import (
	"fmt"
	"strings"
)

// 路由节点
type node struct {
	pattern  string  // 待匹配路由，例如/p/:lang，是从根节点到当前的完整pattern，不是则为空
	part     string  // 路由的最小组成部分，例如：lang，URL切割后的块值
	children []*node // 该节点的子节点，例如[doc,tutorial,intro]
	isWild   bool    // 是否模糊匹配，part含有:或*时候为true
}

// 从子节点找第一个匹配成功的节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 从子节点找所有匹配成功的节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 从parts的height位置寻找/构造节点且连接节点
// 直到完成parts中每个part都有对应的节点
func (n *node) linkNode(pattern string, parts []string, height int) {
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
	child.linkNode(pattern, parts, height+1)
}

// 搜索路径的最终路由节点
func (n *node) searchNode(parts []string, height int) *node {
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
		result := child.searchNode(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}
