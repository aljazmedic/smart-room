package service

import (
	"io/ioutil"
	"time"

	"github.com/AyyItsAljo/smart-room/libraries/go/errors"
	"github.com/AyyItsAljo/smart-room/libraries/go/slog"
	"github.com/AyyItsAljo/smart-room/service.config/domain"
)

// ConfigService handles loading of the config
type ConfigService struct {
	// Config is a pointer to the config object in use by the service
	Config *domain.Config

	// Location is the path to the YAML file holding the config
	Location string
}

// Watch reads the config file every d duration and applies changes
func (s *ConfigService) Watch(d time.Duration) {
	for {
		if _, err := s.Reload(); err != nil {
			slog.Error("Failed to reload config: %v", err)
		}

		time.Sleep(d)
	}
}

// Reload reads the config and applies changes
func (s *ConfigService) Reload() (bool, error) {
	data, err := ioutil.ReadFile(s.Location)
	if err != nil {
		return false, err
	}

	reloaded, err := s.Config.SetFromBytes(data)
	if err != nil {
		return false, errors.InternalService("failed to set config from bytes: %v", err)
	}

	if reloaded {
		slog.Info("Config reloaded")
	}

	return reloaded, nil
}