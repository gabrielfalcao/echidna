package echidna

import (
	"strings"
	"net/http"
	. "launchpad.net/gocheck"
)

// RegexRouter#MakePath should mount the full URL
func (s *S) TestRegexRouterMakePathLocalHostOtherPort(c *C) {
	// Given a router with a port different than 3000
	router := RegexRouter{Domain: "localhost", Port: 3000}
	// When MakePath() is called with a path
	gotten := router.MakePath("/coolf")
	// Then it returns a full URL
	c.Assert(gotten, Equals, "http://localhost:3000/coolf")
}

// RegexRouter#MakePath should mount the full URL inteligently
func (s *S) TestRegexRouterMakePathLocalHost(c *C) {
	// Given a router with the port 80
	router := RegexRouter{Domain: "ffeast.com", Port: 80}
	// When MakePath() is called with a path
	gotten := router.MakePath("/login")
	// Then it returns a full URL
	c.Assert(gotten, Equals, "http://ffeast.com/login")
}



// RegexRouter#Register should register callbacks for a regex + http request method type
func (s *S) TestRegexRouterRegistry(c *C) {
	// Given a router
	rr := RegexRouter{Domain: "localhost", Port: 80}
	// And a fake writer
	fakeWriter := &TestWriter{}

	// When a route is registered
	route := rr.Register("/foo/bar", "GET", func (matches []string, w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("callback called!"))
	})

	// Then its returned route has a pattern
	c.Assert(route.Pattern, Equals, "/foo/bar")

	// And it should respond to the given method has a pattern
	c.Assert(route.Method, Equals, "GET")

	// And its callback was succesfully registered
	route.Callback(fakeWriter, nil)
	c.Assert(fakeWriter.Output, Equals, "callback called!")
}

// RegexRouter#Resolve should resolve a url to a route
func (s *S) TestRegexRouterResolve(c *C) {
	// Given a router
	rr := RegexRouter{Domain: "localhost", Port: 80}
	// And a fake writer
	fakeWriter := &TestWriter{}

	// And a registered route
	rr.Register("/([0-9a-zA-Z]+)", "GET", func (matches []string, w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("callback called!"))
	})
	// When the router can successfully resolve

	route := rr.Resolve("/cool", "GET")

	// Then the callback works
	route.Callback(fakeWriter, nil)
	c.Assert(fakeWriter.Output, Equals, "callback called!")
}


// RegexRouter#Resolve should only match if method is registered
func (s *S) TestRegexRouterResolveFailedByMethod(c *C) {
	// Given a router
	rr := RegexRouter{Domain: "localhost", Port: 80}
	// And a fake writer
	fakeWriter := &TestWriter{}

	// And a registered route
	rr.Register("/([0-9]+)", "GET", func (matches []string, w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("callback not called!"))
	})
	// When the router tries to resolve with a url that is
	// registered but the method is not found
	buf := strings.NewReader("foo=bar")
	request, _ := http.NewRequest("POST", "http://localhost/123", buf)
	route := rr.Resolve("/123", "POST")

	// Then the route should resolve to a 405
	route.Callback(fakeWriter, request)
	c.Assert(fakeWriter.Output, Equals, "<h1>Method Not Allowed</h1>\nPOST /123")
	c.Assert(fakeWriter.Status, Equals, 405)
}


// RegexRouter#Resolve should only match if method is registered
func (s *S) TestRegexRouterResolveFailedNotFound(c *C) {
	// Given a router
	rr := RegexRouter{Domain: "localhost", Port: 80}
	// And a fake writer
	fakeWriter := &TestWriter{}

	// And a registered route
	rr.Register("/([0-9]+)", "GET", func (matches []string, w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("callback not called!"))
	})
	// When the router tries to resolve with a url that is
	// registered but the method is not found
	buf := strings.NewReader("foo=bar")
	request, _ := http.NewRequest("GET", "http://localhost/cool", buf)
	route := rr.Resolve("/cool", "GET")

	// Then the route should resolve to a 404
	route.Callback(fakeWriter, request)
	c.Assert(fakeWriter.Output, Equals,
		"<h1>Not Found</h1>\nGET /cool")

	c.Assert(fakeWriter.Status, Equals, 404)
}
