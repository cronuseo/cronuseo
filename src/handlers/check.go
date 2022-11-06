package handlers

import (
	"errors"
	"fmt"
	"github.com/shashimalcse/Cronuseo/repositories"
	"strings"
)

func CheckAllowed(resource string, role string, action string) (bool, error) {
	var exists bool
	err := repositories.Check(resource, role, action, &exists)
	if err != nil {
		return false, errors.New("not exists")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func Check(resource string, role string, action string) []string {
	allowedActions := []string{}
	actions := strings.Fields(action)
	for _, action := range actions {
		action = strings.TrimSpace(action)
		var exists bool
		err := repositories.Check(resource, role, action, &exists)
		fmt.Println(exists)
		if err == nil {
			if exists {
				allowedActions = append(allowedActions, action)
			}
		}

	}
	return allowedActions
}
