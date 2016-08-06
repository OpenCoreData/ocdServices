package utilities

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/emicklei/go-restful"
)

// This example shows how to create a filter that produces log lines
// according to the Common Log Format, also known as the NCSA standard.
//
// kindly contributed by leehambley
//
// GET http://localhost:8080/ping

var logger *log.Logger = log.New(os.Stdout, "", 0)

func NCSACommonLogFormatLogger() restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		var username = "-"
		if req.Request.URL.User != nil {
			if name := req.Request.URL.User.Username(); name != "" {
				username = name
			}
		}
		chain.ProcessFilter(req, resp)
		logger.Printf("%s - %s [%s] \"%s %s %s\" %d %d \"%s\" \"%s\"",
			strings.Split(req.Request.RemoteAddr, ":")[0],
			username,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			req.Request.Method,
			req.Request.URL.RequestURI(),
			req.Request.Proto,
			resp.StatusCode(),
			resp.ContentLength(),
			req.Request.Referer(),
			req.Request.UserAgent(),
		)
	}
}
