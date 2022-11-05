package handlers

import (
	"errors"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func Check(resource string, role string, action string) (bool, error) {
	var exists bool
	err := repositories.Check(resource, role, action, exists)
	if err != nil {
		return false, errors.New("not exists")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}
