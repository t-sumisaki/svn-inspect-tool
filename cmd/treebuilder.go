package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Node struct {
	Name     string
	Children map[string]*Node
}

func BuildTree(paths []string) *Node {
	root := &Node{
		Name:     "",
		Children: map[string]*Node{},
	}

	for _, p := range paths {
		parts := strings.Split(filepath.Clean(p), string(os.PathSeparator))

		current := root
		for _, part := range parts {
			if part == "" {
				continue
			}

			if _, ok := current.Children[part]; !ok {
				current.Children[part] = &Node{
					Name:     part,
					Children: map[string]*Node{},
				}
			}

			current = current.Children[part]
		}
	}

	return root
}

func PrintTree(n *Node, indent string) string {

	var builder strings.Builder

	var walk func(_node *Node, _indent string, _level int)
	walk = func(_node *Node, _indent string, _level int) {

		keys := make([]string, 0, len(_node.Children))
		for k := range _node.Children {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for i, k := range keys {
			child := _node.Children[k]
			last := i == len(_node.Children)-1
			first := _level == 0
			if first {
				builder.WriteString(child.Name + "\n")
				walk(child, "", _level+1)
			} else if last {
				builder.WriteString(_indent + "└─ " + child.Name + "\n")
				walk(child, _indent+"   ", _level+1)
			} else {
				builder.WriteString(_indent + "├─ " + child.Name + "\n")
				walk(child, _indent+"│  ", _level+1)
			}
		}
	}

	walk(n, indent, 0)

	return builder.String()
}
