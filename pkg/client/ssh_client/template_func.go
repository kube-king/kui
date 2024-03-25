package ssh_client

import (
	"fmt"
	"net/url"
	"strings"
	"text/template"
)

var TemplateFuncMap template.FuncMap

func replace(input, from, to string) string {
	return strings.Replace(input, from, to, -1)
}

func domain(val string) string {
	return getUrlDomain(val)
}

func getUrlDomain(urlAddress string) (host string) {
	parse, err := url.Parse(urlAddress)
	if err != nil {
		host = ""
		return
	}
	if parse.Scheme == "" {
		host = urlAddress
	} else {
		host = parse.Host
	}
	return
}

func lower(val interface{}) string {
	return strings.ToLower(fmt.Sprintf("%v", val))
}

func Default(val interface{}, def interface{}) interface{} {
	if val == nil {
		return def
	}
	return val
}

func init() {
	TemplateFuncMap = template.FuncMap{
		"replace": replace,
		"lower":   lower,
		"domain":  domain,
		"default": Default,
	}
}
