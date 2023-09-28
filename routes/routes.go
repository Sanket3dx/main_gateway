package Routes

import (
	"fmt"
	"main_gateway/controllers"
	"main_gateway/middlewares"
	"main_gateway/utils"

	"github.com/gin-gonic/gin"
)

var routes []utils.Route

func InitRouter() {
	config := utils.GetProjectConfig()
	fmt.Println("👊 Config files are Loaded for Main Geteway ...✅ ")
	Router := gin.Default()
	LoadRoutes(config)
	for _, route := range routes {
		Router.Any(route.Path, middlewares.AuthMiddleware(), controllers.ProxyToService(route))
	}
	Router.POST("/login", controllers.Login)
	fmt.Println("🤘 Routes Loaded ...✅ ")
	fmt.Println("😎 Router started on Port : " + config.Port + " ...✅ ")
	Router.Run(config.Port)
}

func LoadRoutes(config utils.ProjectConfig) {
	for _, service := range config.Services {
		routes = append(routes, utils.Route{
			Path:    service.Service + "/*path",
			Target:  service.URL,
			Methods: []string{"GET", "POST", "PUT", "DELETE"}, // Define allowed methods as needed
		})
	}

}
