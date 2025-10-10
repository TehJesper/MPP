package services

import (
	"fmt"
	"test/domain"
	"test/infrastructure"
)

func checkClass(class string) error {
	classes, err := infrastructure.LoadClassesAndSubclasses()

	if err != nil {
		return err
	}

	if !domain.IsValidClass(class, classes) {
		return fmt.Errorf("invalid class: %s", class)
	}

	return nil
}
