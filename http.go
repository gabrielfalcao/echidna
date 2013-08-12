package echidna

import (
	"fmt"
	"regexp"
	"strings"
	"net/http"
)

type HttpCallback func([]string, http.ResponseWriter, *http.Request)

type Route struct {
	Pattern string;
	regex regexp.Regexp;
	Method string;
	matches []string;
	cb HttpCallback;
}
func (self *Route) Callback (response http.ResponseWriter, request *http.Request) {
	self.cb(self.matches, response, request)
}
func (self *Route) MatchesPath (path string) bool {
	found := self.regex.FindStringSubmatch(path)
	if len(found) > 0 {
		self.matches = found
		return true
	}
	return false
}
func (self *Route) MatchesMethod (method string) bool {
	return self.Method == method;
}
func (self *Route) Matches (path string, method string) bool {
	return self.MatchesPath(path) && self.MatchesMethod(method);
}

type RegexRouter struct {
	Domain string;
	Port int;
	routes []Route;
}
func (self *RegexRouter) MakePath (path string) string {
	fixedPath := strings.TrimLeft(path, "/")

	switch self.Port {
	case 80:
		return fmt.Sprintf("http://%s/%s", self.Domain, fixedPath)
	default:
		return fmt.Sprintf("http://%s:%d/%s", self.Domain, self.Port, fixedPath)
	}
}
func (self *RegexRouter) Register (pattern string, method string, cb HttpCallback) Route {
	regex := regexp.MustCompilePOSIX(pattern)
	r := Route{Pattern: pattern, regex: *regex, Method: method, cb: cb}
	self.routes = append(self.routes, r)
	return r
}
func (self *RegexRouter) Resolve (path string, method string) Route {
	pathMatched := false
	for _, route := range self.routes {
		if route.Matches(path, method) {
			return route
		}
		if route.MatchesPath(path) {
			pathMatched = true
		}
	}
	if pathMatched {
		return self.RouteMethodNotAllowed(path, method)
	} else {
		return self.RouteNotFound(path, method)
	}
}
func (self *RegexRouter) RouteNotFound (path string, method string) Route {
	return Route{Pattern: path, Method: method, cb: func (matches []string, response http.ResponseWriter, request *http.Request) {
			response.WriteHeader(404)
			response.Header().Set("Content-Type", "text-html; charset: utf-8")
			response.Write([]byte("<h1>Not Found</h1>\n"))
			response.Write([]byte(method + " " + path))
	}}
}
func (self *RegexRouter) RouteMethodNotAllowed (path string, method string) Route {
	return Route{Pattern: path, Method: method, cb: func (matches []string, response http.ResponseWriter, request *http.Request) {
			response.WriteHeader(405)
			response.Header().Set("Content-Type", "text-html; charset: utf-8")
			response.Write([]byte("<h1>Method Not Allowed</h1>\n"))
			response.Write([]byte(method + " " + path))
	}}
}
