package services

import (
	"fmt"
	"test/domain"
	cls "test/infrastructure/class"
)

func checkClass(class string) error {
	classes, err := cls.LoadClassesAndSubclasses()

	if err != nil {
		return err
	}

	if !domain.IsValidClass(class, classes) {
		return fmt.Errorf("invalid class: %s", class)
	}

	return nil
}
