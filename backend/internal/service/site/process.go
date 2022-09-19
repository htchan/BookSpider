package site

import (
	"fmt"
	"log"
)

func Process(st *Site) error {
	var err error

	log.Printf("[%v] start backup", st.Name)
	err = Backup(st)
	if err != nil {
		return fmt.Errorf("Process fail: %w", err)
	}

	log.Printf("[%v] start validate", st.Name)
	err = Validate(st)
	if err != nil {
		return fmt.Errorf("Process fail: %w", err)
	}

	log.Printf("[%v] start update", st.Name)
	err = Update(st)
	if err != nil {
		return fmt.Errorf("Process fail: %w", err)
	}

	log.Printf("[%v] start explore", st.Name)
	err = Explore(st)
	if err != nil {
		return fmt.Errorf("Process fail: %w", err)
	}

	log.Printf("[%v] start check", st.Name)
	err = Check(st)
	if err != nil {
		return fmt.Errorf("Process fail: %w", err)
	}

	// log.Printf("[%v] start download", st.Name)
	// err = Download(st)
	// if err != nil {
	// 	return fmt.Errorf("Process fail: %w", err)
	// }

	log.Printf("[%v] start fix", st.Name)
	err = Fix(st)
	if err != nil {
		return fmt.Errorf("Process fail: %w", err)
	}

	return nil
}
