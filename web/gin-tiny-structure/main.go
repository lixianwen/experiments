package main

import (
	"fmt"
	"gdemo/internal/config"
	"gdemo/internal/router"
	"log/slog"
	"os"
)

func main() {
	config := config.GetConfig()
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: config.Logger.AddSource,
		Level:     slog.Level(config.Logger.Level),
	})
	slog.SetDefault(slog.New(h))

	r := router.SetupRouter()
	r.Run(fmt.Sprintf("%s:%d", config.HTTP.Host, config.HTTP.Port))
}
