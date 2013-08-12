package echidna

import (
	"testing"
	"net/http"

	. "launchpad.net/gocheck"
)

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})

type TestWriter struct {
	Output string
	header http.Header
	Status int
}

func (self *TestWriter) Header() http.Header {
	if self.header == nil {
		self.header = make(http.Header)
	}
	return self.header;
}

func (self *TestWriter) Write(content []byte) (int, error) {
	self.Output += string(content)
	return len(self.Output), nil
}

func (self *TestWriter) WriteHeader(status int) {
	self.Status = status
}
