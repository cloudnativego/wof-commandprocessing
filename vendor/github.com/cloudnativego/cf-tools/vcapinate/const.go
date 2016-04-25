package main

import (
	"fmt"
	"text/template"
)

const (
	serviceTemplate = `
{
  "user-provided": [
  {{ range $i, $e := .UserProvided }}
  {{- if $i }}
    },
  {{ end -}}
    {
    "credentials": {
		{{ print (convertCredentialsToString $e.Credentials) }}
    },
    "label": "user-provided",
    "name": "{{ $e.Name }}",
    "syslog_drain_url": "{{ $e.SyslogURL }}",
    "tags": []
  {{ end -}}
    }
  ]
}
`
)

var fns = template.FuncMap{
	"convertCredentialsToString": func(c map[interface{}]interface{}) (out string) {
		l := len(c)
		i := 0
		out = ""
		for k, v := range c {
			out += fmt.Sprintf("%q:%q", k.(string), v.(string))
			i++
			if i < l {
				out += ", "
			}
		}
		return
	},
}
