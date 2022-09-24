package site

import (
	"fmt"
	"log"
)

func Process(st *Site) error {
	var err error

	log.Printf("[operation.%v.backup] start", st.Name)
	err = Backup(st)
	log.Printf("[operation.%v.backup] complete", st.Name)
	if err != nil {
		return fmt.Errorf("Process fail: %w", err)
	}

	log.Printf("[operation.%v.validate] start", st.Name)
	err = Validate(st)
	log.Printf("[operation.%v.validate] complete", st.Name)
	if err != nil {
		return fmt.Errorf("Process fail: %w", err)
	}

	log.Printf("[operation.%v.update] start", st.Name)
	err = Update(st)
	log.Printf("[operation.%v.update] complete", st.Name)
	if err != nil {
		return fmt.Errorf("Process fail: %w", err)
	}

	log.Printf("[operation.%v.explore] start", st.Name)
	err = Explore(st)
	log.Printf("[operation.%v.explore] complete", st.Name)
	if err != nil {
		return fmt.Errorf("Process fail: %w", err)
	}

	log.Printf("[operation.%v.check] start", st.Name)
	err = Check(st)
	log.Printf("[operation.%v.check] complete", st.Name)
	if err != nil {
		return fmt.Errorf("Process fail: %w", err)
	}

	log.Printf("[operation.%v.download] start", st.Name)
	err = Download(st)
	log.Printf("[operation.%v.download] complete", st.Name)
	if err != nil {
		return fmt.Errorf("Process fail: %w", err)
	}

	log.Printf("[operation.%v.fix] start", st.Name)
	err = Fix(st)
	log.Printf("[operation.%v.fix] complete", st.Name)
	if err != nil {
		return fmt.Errorf("Process fail: %w", err)
	}

	return nil
}
