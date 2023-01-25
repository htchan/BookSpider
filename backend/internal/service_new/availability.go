package service

import (
	"fmt"
	"strings"
)

func (serv *ServiceImp) CheckAvailability() error {
	html, err := serv.client.Get(serv.conf.AvailabilityConfig.URL)
	if err != nil {
		return fmt.Errorf("validate fail: %w", err)
	}

	if !strings.Contains(html, serv.conf.AvailabilityConfig.CheckString) {
		return ErrInvalidSite
	}

	return nil
}
