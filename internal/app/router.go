package app

import "github.com/shubham-gaur/goui/internal/controllers"

func mapURLs() {
	router.GET("/", controllers.Navigator.Index)
	router.POST("/server", controllers.Navigator.Server)
	router.POST("/kill/:serverid", controllers.Navigator.KillServer)
}
