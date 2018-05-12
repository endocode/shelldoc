package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Endocode/shelldoc/pkg/shell"

	"gopkg.in/russross/blackfriday.v2"
)

func visitor(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	fmt.Printf("node: %s - %s - entering: %v.\n", node.Type, node.Literal, entering)
	if node.Type == blackfriday.CodeBlock {
		// fmt.Printf("node: %s - %s - entering: %v.\n", node.Type, node.Literal, entering)
	}
	return blackfriday.GoToNext
}

func main() {
	if data, err := ReadInput(os.Args[1:]); err != nil {
		log.Fatalln(err)
	} else {
		shell, err := shell.StartShell()
		if err != nil {
			log.Fatalln(err)
		}
		defer shell.Exit()

		md := blackfriday.New()
		om := md.Parse(data)
		om.Walk(visitor)
		fmt.Println("--------------------")
		fmt.Println(string(data))
		fmt.Println("Done.")
	}
}
