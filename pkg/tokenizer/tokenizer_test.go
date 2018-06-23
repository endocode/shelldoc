package tokenizer

// This file is part of shelldoc.
// © 2018, Mirko Boehm <mirko@endocode.com> and the shelldoc contributors
// SPDX-License-Identifier: Apache-2.0

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

var echoTrueCodeBlockCount int

func codeBlockHandler(visitor *Visitor, node *blackfriday.Node) blackfriday.WalkStatus {
	//fmt.Printf("%s: %v\n", node.Type, string(node.Literal))
	echoTrueCodeBlockCount++
	return blackfriday.GoToNext
}
func TestEchoTrue(t *testing.T) {
	data, err := ioutil.ReadFile("samples/echotrue.md")
	require.NoError(t, err, "Unable to read sample data file")
	visitor := Visitor{codeBlockHandler, codeBlockHandler, nil}
	require.Zero(t, echoTrueCodeBlockCount, "Starting the counter")
	Tokenize(data, &visitor)
	require.Equal(t, echoTrueCodeBlockCount, 1, "There is one code block element in the sample file")
}

func TestTokenizeEchoTrue(t *testing.T) {
	data, err := ioutil.ReadFile("samples/echotrue.md")
	require.NoError(t, err, "Unable to read sample data file")
	visitor := NewInteractionVisitor()
	Tokenize(data, visitor)
	require.Equal(t, len(visitor.Interactions), 1, "There is one code block element in the sample file")
}

func TestTokenizeHelloWorld(t *testing.T) {
	data, err := ioutil.ReadFile("samples/helloworld.md")
	require.NoError(t, err, "Unable to read sample data file")
	visitor := NewInteractionVisitor()
	Tokenize(data, visitor)
	require.Equal(t, 4, len(visitor.Interactions), "There are three code block elements with a total of 4 interactions in the sample file")
	require.Empty(t, visitor.Interactions[0].Response, "The first command does not expect a response")
	require.NotEmpty(t, visitor.Interactions[1].Response, "The second command expects a response")
	require.Equal(t, visitor.Interactions[1].Response[0], "Hello", "The second command expects a response")
	require.NotEmpty(t, visitor.Interactions[2].Response, "The third command expects a response")
	require.Equal(t, visitor.Interactions[2].Response[0], "World", "The third command expects a response")
	fourth := visitor.Interactions[3]
	require.Equal(t, 2, len(fourth.Response), "The response of the fourth interaction contains two lines")
	require.Equal(t, "...", fourth.Response[1], "The last line of the fourth response is an ellipsis")
}

func TestTokenizeFenced(t *testing.T) {
	data, err := ioutil.ReadFile("samples/fenced.md")
	require.NoError(t, err, "Unable to read sample data file")
	visitor := NewInteractionVisitor()
	Tokenize(data, visitor)
	require.Equal(t, len(visitor.Interactions), 2, "There are two fenced code block in the sample file.")
	first := visitor.Interactions[0]
	require.Equal(t, first.Language, "shell", "shell was the specified languagwe for the first code block")
	require.Equal(t, first.Attributes["shelldocexitcode"], "1", "1 was the specified value for shelldocexitcode")
	require.Equal(t, first.Attributes["shelldocwhatever"], "", "shelldocwhatever comes with no value")
	_, exists := first.Attributes["shelldocnonsense"]
	require.False(t, exists, "shelldocnonsense was not defined")
	second := visitor.Interactions[1]
	require.Empty(t, second.Language, "No language was specified in the second block")
	require.Empty(t, second.Attributes, "No attributes where specified in the second block")
}
