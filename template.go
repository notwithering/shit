package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var (
	curlReg = regexp.MustCompile("curl/?")
)

type tmplData struct {
	Links  []string
	Upload bool
}

func execute(w http.ResponseWriter, r *http.Request, links []string) error {
	response := getResponse(r, links)
	_, err := w.Write([]byte(response))
	return err
}

func getResponse(r *http.Request, links []string) string {
	if curlReg.MatchString(r.UserAgent()) {
		return makeCurlResponse(links)
	}
	return makeHtmlResponse(links)
}

func makeHtmlResponse(links []string) string {
	return `<!doctype html>
<html data-fbscriptallow="true">
		<head>
			<meta name="viewport" content="width=device-width" />
		</head>
		<body>
			<pre>` + func() string {
		b := strings.Builder{}
		for _, name := range links {
			fmt.Fprintf(&b, "<a href=\"%s\">%s</a>\n", name, name)
		}
		return b.String()
	}() + `</pre>

			` + func() string {
		if upload {
			return `		<form method="POST" enctype="multipart/form-data">
			<input type="file" name="files" multiple required />
			<button type="submit">Upload</button>
		</form>
`
		}
		return ""
	}() + `	</body>
</html>
`
}

func makeCurlResponse(links []string) string {
	b := strings.Builder{}
	for _, name := range links {
		b.WriteString(name)
		b.WriteByte('\n')
	}
	return b.String()
}
