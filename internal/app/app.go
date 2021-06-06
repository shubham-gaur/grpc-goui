package app

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/shubham-gaur/goui/internal/services/templates"
)

var (
	router = gin.Default()
	server = flag.String("server", ":8080", "Accepted server listen IP :<port> or <ip>:<port>")
)

// StartApplication ...
func StartApplication() {
	flag.Parse()
	router.HTMLRender = templates.UserTemplate.CreateUserTemplate()
	mapURLs()
	router.Run(*server)
}
