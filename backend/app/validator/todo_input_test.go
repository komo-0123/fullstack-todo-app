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
		errMsg    string
		expectErr bool
	}{
		"エラーなし":      {model.Todo{ID: 1, Title: "タイトル", IsComplete: false}, "", noErr},
		"タイトルが空":     {model.Todo{ID: 1, Title: "", IsComplete: false}, "タイトルは必須です。", wantErr},
		"タイトルが256文字": {model.Todo{ID: 1, Title: strings.Repeat("あ", 256), IsComplete: false}, "タイトルは255文字以内で入力してください。", wantErr},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := validator.TodoInput(c.input)
			if c.expectErr {
				if err == nil || err.Error() != c.errMsg {
					t.Errorf("want: %s, got: %s", c.errMsg, err.Error())
				}
			} else if err != nil {
				t.Errorf("want: nil, got: %s", err.Error())
			}
		})
	}
}
