package services
import (
	"test/domain"
	"test/infrastructure"
	"fmt"
)

func checkClass(class string) (error) {
	classes, err := infrastructure.LoadClassesAndSubclasses()

    if err != nil {
        return err
    }

    if !domain.IsValidClass(class, classes) {
        return fmt.Errorf("invalid class: %s", class)
    }

	return nil
}
