// Code generated by goagen v1.3.1, DO NOT EDIT.
//
// cellar JavaScript Client Example
//
// Command:
// $ goagen
// --design=github.com/goadesign/goa-cellar/design
// --out=$(GOPATH)/src/github.com/goadesign/goa-cellar/public
// --version=v1.3.1

package js

import (
	"github.com/goadesign/goa"
)

// MountController mounts the JavaScript example controller under "/js".
// This is just an example, not the best way to do this. A better way would be to specify a file
// server using the Files DSL in the design.
// Use --noexample to prevent this file from being generated.
func MountController(service *goa.Service) {
	// Serve static files under js
	service.ServeFiles("/js/*filepath", "/home/raphael/go/src/github.com/goadesign/goa-cellar/public/js")
	service.LogInfo("mount", "ctrl", "JS", "action", "ServeFiles", "route", "GET /js/*")
}
