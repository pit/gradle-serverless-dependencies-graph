package main

var Template = `
<html><body><pre>
{{range .Items}}
<a href="/dependency/{{.Child}}">{{.Child}}</a>
{{end}}
</pre></body></html>
`
