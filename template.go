package mommy

import (
	"io"
	"regexp"
	"strings"
)

var templateRegex = regexp.MustCompile(`{([a-zA-Z0-9_]+)}`)

type templateSpan struct {
	key  VariableKey
	span [2]int
}

type templater struct {
	text  Template
	spans []templateSpan
}

func compileTemplate(t Template) templater {
	matches := templateRegex.FindAllStringSubmatchIndex(string(t), -1)
	spans := make([]templateSpan, len(matches))
	for i, m := range matches {
		spans[i] = templateSpan{
			key:  VariableKey(t[m[2]:m[3]]),
			span: [2]int{m[0], m[1]},
		}
	}
	return templater{
		text:  t,
		spans: spans,
	}
}

func (t templater) render(vars map[VariableKey]string) string {
	var s strings.Builder
	s.Grow(len(t.text))
	t.renderTo(&s, vars)
	return s.String()
}

func (t templater) renderTo(w io.Writer, vars map[VariableKey]string) error {
	var last int
	var err error
	for _, span := range t.spans {
		_, err = io.WriteString(w, string(t.text[last:span.span[0]]))
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, vars[span.key])
		if err != nil {
			return err
		}
		last = span.span[1]
	}
	_, err = io.WriteString(w, string(t.text[last:]))
	return err
}
