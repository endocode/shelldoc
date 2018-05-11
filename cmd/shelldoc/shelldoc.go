package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/russross/blackfriday.v2"
)

func visitor(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	if node.Type == blackfriday.CodeBlock {
		fmt.Printf("node: %s - %s - entering: %v.\n", node.Type, node.Literal, entering)
	}
	return blackfriday.GoToNext
}

func main() {
	data, err := ioutil.ReadFile("../../README.md")
	if err != nil {
		log.Fatalln(err)
	}

	md := blackfriday.New()
	om := md.Parse(data)
	om.Walk(visitor)
	fmt.Println("Done.")
}
