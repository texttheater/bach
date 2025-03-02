package docutil

import (
	"fmt"
	"html"
	"io"
	"strings"

	"github.com/texttheater/bach/shapes"
)

// InlineCode takes a string representing some program code and converts it to
// a Markdown representation suitable for processing by mdbook.
func InlineCode(s string) string {
	// handle empty string specially: do not generate <code></code> tags
	if s == "" {
		return ""
	}
	// escape HTML special characters
	s = html.EscapeString(s)
	// escape characters that mdbook treats specially
	s = strings.ReplaceAll(s, "|", "&#124;")
	s = strings.ReplaceAll(s, "\\", "&#92;")
	// wrap in code tags
	s = "<code>" + s + "</code>"
	// return
	return s
}

func PrintExamplesTable(w io.Writer, examples []shapes.Example) {
	fmt.Fprintf(w, "| Program | Type | Value | Error |\n")
	fmt.Fprintf(w, "|---|---|---|---|\n")
	for _, example := range examples {
		var err string
		if example.Error == nil {
			err = ""
		} else {
			err = strings.TrimSpace(fmt.Sprintf("%s", example.Error))
		}
		fmt.Fprintf(
			w,
			"| %s | %s | %s | %s |\n",
			InlineCode(example.Program),
			InlineCode(example.OutputType),
			InlineCode(example.OutputValue),
			InlineCode(err),
		)
	}
}
