package main

import (
	"fmt"
	"gdemo/internal/config"
	"gdemo/internal/router"
	"log/slog"
	"os"

	"gdemo/docs"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

//	@title			Gin example
//	@version		0.1
//	@description	This is a sample gin server.

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Bearer authentication
func main() {
	config := config.GetConfig()
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: config.Logger.AddSource,
		Level:     slog.Level(config.Logger.Level),
	})
	slog.SetDefault(slog.New(h))

	addr := fmt.Sprintf("%s:%d", config.HTTP.Host, config.HTTP.Port)

	docs.SwaggerInfo.Host = addr
	docs.SwaggerInfo.Schemes = []string{"http"}

	r := router.SetupRouter()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(addr)
}
