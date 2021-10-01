package main

var Template = `
<html><body><pre>
{{range items}}
<a href="/repository/{{.Child}}">{{.Child}}</a>
{{end}
</pre></body></html>
`
