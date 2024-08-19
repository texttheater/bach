package main

import (
	"fmt"
	"html"
	"log"
	"os"
	"strings"

	"github.com/texttheater/bach/builtin"
	"github.com/texttheater/bach/shapes"
	"github.com/urfave/cli/v2"
)

var FuncersByCategory = map[string][]shapes.Funcer{
	"null":    builtin.NullFuncers,
	"io":      builtin.IOFuncers,
	"logic":   builtin.LogicFuncers,
	"math":    builtin.MathFuncers,
	"text":    builtin.TextFuncers,
	"arr":     builtin.ArrFuncers,
	"obj":     builtin.ObjFuncers,
	"types":   builtin.TypeFuncers,
	"values":  builtin.ValueFuncers,
	"regexp":  builtin.RegexpFuncers,
	"control": builtin.ControlFuncers,
}

// inlineCode takes a string representing some program code and converts it to
// a Markdown representation suitable for processing by mdbook.
func inlineCode(s string) string {
	// handle empty string specially: do not generate <code></code> tags
	if s == "" {
		return ""
	}
	// escape HTML entities
	s = html.EscapeString(s)
	// escape characters that mdbook treats specially
	s = strings.ReplaceAll(s, "|", "&#124;")
	s = strings.ReplaceAll(s, "\\", "&#92;")
	// wrap in code tags
	s = "<code>" + s + "</code>"
	// return
	return s
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "builtin",
				Usage: "generate documentation for the builtin funcers of a given category",
				Action: func(cCtx *cli.Context) error {
					cat := cCtx.Args().First()
					funcers, ok := FuncersByCategory[cat]
					for _, funcer := range funcers {
						if !ok {
							return cli.Exit("unknown category", 1)
						}
						fmt.Printf("## %s\n\n", funcer.Name)
						fmt.Printf("%s\n\n", funcer.Summary)
						fmt.Printf("| | Type | Value |\n")
						fmt.Printf("|---|---|---|\n")
						fmt.Printf("| Input | %s | %s |\n", inlineCode(funcer.InputType.String()), funcer.InputDescription)
						for i, param := range funcer.Params {
							fmt.Printf("| %s (param #%d) | %s | %s |\n", param.Name, i+1, inlineCode(param.String()), param.Description)
						}
						fmt.Printf("|Output | %s | %s |\n\n", inlineCode(funcer.OutputType.String()), funcer.OutputDescription)
						fmt.Printf("%s\n\n", funcer.Notes)
						fmt.Printf("### Examples\n\n")
						fmt.Printf("| Program | Type | Value | Error |\n")
						fmt.Printf("|---|---|---|---|\n")
						for _, example := range funcer.Examples {
							var err string
							if example.Error == nil {
								err = ""
							} else {
								err = strings.TrimSpace(fmt.Sprintf("%s", example.Error))
							}
							fmt.Printf(
								"| %s | %s | %s | %s |\n",
								inlineCode(example.Program),
								inlineCode(example.OutputType),
								inlineCode(example.OutputValue),
								inlineCode(err),
							)
						}
						fmt.Printf("\n")
					}
					return nil

				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
