package main

import (
	"github.com/AyyItsAljo/smart-room/libraries/go/bootstrap"
	"github.com/AyyItsAljo/smart-room/libraries/go/config"
	"github.com/AyyItsAljo/smart-room/libraries/go/router"
	"github.com/AyyItsAljo/smart-room/libraries/go/slog"

	"github.com/AyyItsAljo/smart-room/service.config/domain"
	"github.com/AyyItsAljo/smart-room/service.config/handler"
	"github.com/AyyItsAljo/smart-room/service.config/service"

)

func main(){
	c := domain.Config{}

	configService := service.ConfigService{
		Location: "/data/config.yaml",
		Config:   &c,
	}

	_, err := configService.Reload()
	if err != nil {
		slog.Panic("Failed to load config: %v", err)
	}

	selfConfig, err := c.Get("service.config")
	if err != nil {
		slog.Panic("Error reading own config: %v", err)
	}

	config.DefaultProvider = config.New(selfConfig)

	if config.Get("polling.enabled").Bool(false) {
		interval := config.Get("polling.interval").Int(30000)
		slog.Info("Polling for config changes every %d milliseconds", interval)
		go configService.Watch(time.Millisecond * time.Duration(interval))
	}

	configHandler := handler.ConfigHandler{
		Config:        &c,
		ConfigService: &configService,
	}

	r := router.New()
	r.Get("/read/{service_name}", configHandler.ReadConfig)
	r.Patch("/reload", configHandler.ReloadConfig)

	bootstrap.Run(r)
}