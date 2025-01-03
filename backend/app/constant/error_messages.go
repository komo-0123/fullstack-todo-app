package constant

// 入力関連のエラーメッセージ
const (
	INPUT_ERR_INVALID_INPUT = "入力が不正です。"
	INPUT_ERR_INVALID_ID    = "IDが不正です。"
	INPUT_ERR_FAILED_GET_ID = "IDの取得に失敗しました。"
)

// DB操作関連のエラーメッセージ
const (
	DB_ERR_FAILED_GET_TODO     = "TODOの取得に失敗しました。"
	DB_ERR_FAILED_GET_TODO_ROW = "TODOの読み込みに失敗しました。"
	DB_ERR_NOT_FOUND_TODO      = "TODOが見つかりません。"
	DB_ERR_FAILED_ADD_TODO     = "TODOの追加に失敗しました。"
	DB_ERR_FAILED_UPDATE_TODO  = "TODOの更新に失敗しました。"
	DB_ERR_NOT_UPDATED_TODO    = "更新したTODOがありません。"
	DB_ERR_FAILED_DELETE_TODO  = "TODOの削除に失敗しました。"
	DB_ERR_DELETED_TODO        = "指定のTODOは削除済みです。"
)
