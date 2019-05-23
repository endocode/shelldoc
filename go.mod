module github.com/endocode/shelldoc

go 1.12

require (
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/spf13/cobra v0.0.4
	github.com/stretchr/testify v1.3.0
	gopkg.in/russross/blackfriday.v2 v2.0.1
)

replace gopkg.in/russross/blackfriday.v2 v2.0.1 => github.com/russross/blackfriday/v2 v2.0.1
