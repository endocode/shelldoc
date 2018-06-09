package main

import (
	"log"
	"os"

	"github.com/Endocode/shelldoc/pkg/shell"
	"github.com/Endocode/shelldoc/pkg/tokenizer"
)

// func visitor(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
// 	fmt.Printf("node: %s - %s - entering: %v.\n", node.Type, node.Literal, entering)
// 	if node.Type == blackfriday.CodeBlock {
// 		// fmt.Printf("node: %s - %s - entering: %v.\n", node.Type, node.Literal, entering)
// 	}
// 	return blackfriday.GoToNext
// }

func main() {
	data, err := ReadInput(os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}
	// run the input through the tokenizer
	visitor := tokenizer.NewInteractionVisitor()
	tokenizer.Tokenize(data, visitor)
	// start a background shell, it will run until the program ends
	shell, err := shell.StartShell()
	if err != nil {
		log.Fatalln(err)
	}
	defer shell.Exit()
	// execute the interactions and verify the results:
	for index, interaction := range visitor.Interactions {
		log.Printf("[%2d]: %s\n", index, interaction.Describe())
		log.Printf("--> : %s", interaction.Cmd)
		if err := interaction.Execute(&shell); err != nil {
			log.Printf("--    ERROR: %v", err)
		}
		log.Printf("<-- : %s", interaction.Result())
	}
}
