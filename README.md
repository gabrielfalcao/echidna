# Echidna

Web tools for Go


## HTTP Routing

```go
import (
	"github.com/gabrielfalcao/echidna"
)
view := echidna.views.TemplateDirectory("templates")

router := echidna.RegexRouter{Domain: "localhost", Port: 3000}
router.Static("/static/", "./assets")

router.Register("/person/([0-9]+)$", "GET", func (matches []string, response http.ResponseWriter, request *http.Request){
        context := make(map[string][interface{}])

        user_name := matches[0]
        context["user"] =
        view.Render("index.html", context)
})


echidna.Activate(router)
```
