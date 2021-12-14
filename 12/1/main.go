package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Node struct {
	name        string
	isBigCave   bool
	connections []*Node
}

type Graph struct {
	nodes []*Node
}

func (g Graph) String() string {
	str := ""
	for _, node := range g.nodes {
		str += fmt.Sprintf("%p: %#v\n", node, *node)
	}
	return str
}

func makeConnection(g *Graph, a, b string) {
	aFound := false
	bFound := false
	var aNode *Node
	var bNode *Node
	for _, node := range g.nodes {
		if node.name == a {
			aFound = true
			aNode = node
		} else if node.name == b {
			bFound = true
			bNode = node
		}
	}
	if !aFound {
		aNode = &Node{a, strings.ToUpper(a) == a, nil}
		g.nodes = append(g.nodes, aNode)
	}
	if !bFound {
		bNode = &Node{b, strings.ToUpper(b) == b, nil}
		g.nodes = append(g.nodes, bNode)
	}

	aNode.connections = append(aNode.connections, bNode)
	bNode.connections = append(bNode.connections, aNode)
}

type Path []*Node

func inBanList(n *Node, banList []*Node) bool {
	for _, node := range banList {
		if n == node {
			return true
		}
	}
	return false
}

func findPathsFromNode(cur *Node, curPath, banList []*Node) []Path {
	if !cur.isBigCave {
		banList = append(banList, cur)
	}

	curPath = append(curPath, cur)
	if cur.name == "end" {
		return []Path{curPath}
	}

	paths := make([]Path, 0)
	for _, con := range cur.connections {
		if inBanList(con, banList) {
			continue
		}
		paths = append(paths,
			findPathsFromNode(con, curPath, banList)...)
	}
	return paths
}

func findPaths(g *Graph) []Path {
	for _, node := range g.nodes {
		if node.name == "start" {
			return findPathsFromNode(node, nil, nil)
		}
	}
	return nil //Unreachable
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	input = input[:len(input)-1] //Remove last

	graph := &Graph{make([]*Node, 0)}
	for _, conn := range strings.Split(string(input), "\n") {
		nodes := strings.Split(conn, "-")
		makeConnection(graph, nodes[0], nodes[1])
	}

	paths := findPaths(graph)
	fmt.Println(len(paths))
}
