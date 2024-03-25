package task

import (
	"bytes"
	"text/template"
)

func TextRender(env interface{}, tmpl string) (result string, err error) {

	t, err := template.New("textTemplate").Parse(tmpl)
	if err != nil {
		return
	}

	buff := new(bytes.Buffer)
	err = t.Execute(buff, env)
	result = buff.String()
	return
}
