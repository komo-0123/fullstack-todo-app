package validator_test

import (
	"backend/app/model"
	"backend/app/validator"
	"strings"

	"testing"
)

func TestTodoInput(t *testing.T) {
	wantErr, noErr := true, false
	cases := map[string]struct {
		input     model.Todo
		expectErr bool
	}{
		"エラーなし":      {model.Todo{ID: 1, Title: "タイトル", IsComplete: false}, noErr},
		"タイトルが空":     {model.Todo{ID: 1, Title: "", IsComplete: false}, wantErr},
		"タイトルが256文字": {model.Todo{ID: 1, Title: strings.Repeat("あ", 256), IsComplete: false}, wantErr},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := validator.TodoInput(c.input)
			if c.expectErr && err == nil {
				t.Errorf("エラーが発生しませんでした。")
			}
			if !c.expectErr && err != nil {
				t.Errorf("エラーが発生しました: %v", err)
			}
		})
	}
}
