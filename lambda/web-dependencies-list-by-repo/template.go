package main

var Template = `
<html><body><pre>
{{.Repo}}/{{.Ref}}:
{{range .Items}}
{{.Dependency}}:{{.Version}}
{{end}}
</pre></body></html>
`
