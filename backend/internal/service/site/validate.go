package site

import (
	"errors"
	"fmt"
	"strings"
)

func Validate(st *Site) error {
	html, err := st.Client.Get(st.StConf.AvailabilityConfig.URL)
	if err != nil {
		return fmt.Errorf("validate fail: %w", err)
	}

	if !strings.Contains(html, st.StConf.AvailabilityConfig.Check) {
		return errors.New("invalid site")
	}

	return nil
}
