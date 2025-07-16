package main

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"strings"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/shapes"
)

// inlineCode takes a string representing some program code and converts it to
// a Markdown representation suitable for processing by mdbook.
func inlineCode(s string) string {
	// handle empty string specially: do not generate <code></code> tags
	if s == "" {
		return ""
	}
	// escape HTML special characters
	s = html.EscapeString(s)
	// turn newlines into <br> tags
	s = strings.ReplaceAll(s, "\n", "<br>")
	// escape characters that mdbook treats specially
	s = strings.ReplaceAll(s, "|", "&#124;")
	s = strings.ReplaceAll(s, "\\", "&#92;")
	s = strings.ReplaceAll(s, "~", "&#126;")
	// wrap in code tags
	s = "<code>" + s + "</code>"
	return s
}

func printExamplesTable(w io.Writer, examples []shapes.Example) {
	fmt.Fprintf(w, "| Program | Type | Value | Error |\n")
	fmt.Fprintf(w, "|---|---|---|---|\n")
	for _, example := range examples {
		var err string
		if example.Error == nil {
			err = ""
		} else {
			var buf bytes.Buffer
			errors.Explain(&buf, example.Error, example.Program)
			err = buf.String()
		}
		fmt.Fprintf(
			w,
			"| %s | %s | %s | %s |\n",
			inlineCode(example.Program),
			inlineCode(example.OutputType),
			inlineCode(example.OutputValue),
			inlineCode(err),
		)
	}
}
