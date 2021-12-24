/*
在本文件中实现的动态路由具备以下两个功能：
第一，参数匹配：例如/p/:lang/doc可以匹配/p/c/doc和/p/go.doc
第二，通配*，例如/static/*filepath可以匹配/static/fav.ico，也可以匹配/static/js/jQuery.js，这种模式常用于
静态服务器，能够递归地匹配子路径
*/

package gee

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 字节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true 表示不精确匹配
}

// 格式化函数
func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// 遍历Trie树
func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 插入节点
// /p/:lang/doc只有在第三层节点，即doc节点，pattern才会设置为/p/:lang/doc。p和:lang节点的pattern属性
// 皆为空。因此，当匹配结束时，我们可以使用n.pattern == ""来判断路由规则是否匹配成功。
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern // 找到了路由该插入的层，插入节点，设置 pattern
		return
	}

	part := parts[height]       // 路由中的一部分
	child := n.matchChild(part) // 找到第一个匹配上的位置
	// 如果之前不存在的话就要新建一个中间节点，否则就无须操作
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1) // 递归处理下一层的信息
}

// 查找节点
// 查询功能，同样也是递归查询每一层的节点，退出规则是，匹配到了*，匹配失败，或者匹配到了第len(parts)层节点。
func (n *node) search(parts []string, height int) *node {
	// 如果到了最下面一层，或者出现了通配符
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 当前节点的 pattern 为空，说明没有到达前缀树的叶子，返回 nil
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height] // 处理第height位置的路由部分
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1) // 递归处理后面的路由部分
		if result != nil {
			return result
		}
	}

	return nil
}
