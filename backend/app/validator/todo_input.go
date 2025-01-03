package validator

import (
	consts "backend/app/constant"
	"backend/app/model"
	"fmt"
	"strings"
)

func TodoInput(todo model.Todo) error {
	if len(strings.TrimSpace(todo.Title)) == 0 {
		return fmt.Errorf(consts.INPUT_ERR_REQUIRED_TITLE)
	}
	if len(todo.Title) > 255 {
		return fmt.Errorf(consts.INPUT_ERR_OVER_LENGTH_TITLE)
	}
	return nil
}
