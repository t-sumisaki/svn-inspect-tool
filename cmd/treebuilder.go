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

	var walk func(_node *Node, _indent string)
	walk = func(_node *Node, _indent string) {

		keys := make([]string, 0, len(_node.Children))
		for k := range _node.Children {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for i, k := range keys {
			child := _node.Children[k]
			last := i == len(_node.Children)-1
			first := _indent == ""
			if first {
				builder.WriteString(child.Name + "\n")
				walk(child, "    ")
			} else if last {
				builder.WriteString(_indent + "└── " + child.Name + "\n")
				walk(child, _indent+"    ")
			} else {
				builder.WriteString(_indent + "├── " + child.Name + "\n")
				walk(child, _indent+"│   ")
			}
		}
	}

	walk(n, indent)

	return builder.String()
}
