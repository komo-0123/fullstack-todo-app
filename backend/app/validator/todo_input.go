package validator

import (
	"backend/app/model"
	"fmt"
	"strings"
)

func TodoInput(todo model.Todo) error {
	const (
		errRequiredTitle   = "タイトルは必須です。"
		errOverLengthTitle = "タイトルは255文字以内で入力してください。"
	)

	if len(strings.TrimSpace(todo.Title)) == 0 {
		return fmt.Errorf(errRequiredTitle)
	}
	if len(todo.Title) > 255 {
		return fmt.Errorf(errOverLengthTitle)
	}

	return nil
}
