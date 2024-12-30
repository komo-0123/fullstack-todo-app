package validator

import (
	"backend/app/model"
	"fmt"
	"strings"
)

func TodoInput(todo model.Todo) error {
	if len(strings.TrimSpace(todo.Title)) == 0 {
		return fmt.Errorf("タイトルは必須です。")
	}
	if len(todo.Title) > 255 {
		return fmt.Errorf("タイトルは255文字以内で入力してください。")
	}
	return nil
}
