package main

var Template = `
<html><body><pre>
{{if $.Parent}}<a href="/repository">../</a>{{end}}
{{range .Items}}
<a href="/repository/{{if $.Parent}}{{$.Parent}}/{{end}}{{.Child}}">{{.Child}}</a>
{{else}}
No items found
{{end}}
</pre></body></html>
`
