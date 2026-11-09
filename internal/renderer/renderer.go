package renderer

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

var defaultFuncMap = template.FuncMap{
	"lower": strings.ToLower,
	"upper": strings.ToUpper,
	"title": strings.Title,
}

// RenderString parses and executes a template string with the provided data.
func RenderString(name, tmplStr string, data any) ([]byte, error) {
	tmpl, err := template.New(name).Funcs(defaultFuncMap).Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	return buf.Bytes(), nil
}
