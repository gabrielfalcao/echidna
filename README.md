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

        context["user"] = matches[0]
        view.Render("index.html", context)
})


echidna.Activate(router)
```


[![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/gabrielfalcao/echidna/trend.png)](https://bitdeli.com/free "Bitdeli Badge")

