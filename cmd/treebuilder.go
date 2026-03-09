package main

import (
	"os"
	"path/filepath"
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
		for _, child := range _node.Children {
			builder.WriteString(_indent + child.Name)
			walk(child, _indent+"\t")
		}
	}

	walk(n, indent)

	return builder.String()
}
