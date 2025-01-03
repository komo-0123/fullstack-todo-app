package constant

// リクエスト関連のエラーメッセージ
const (
	HTTP_ERR_NOT_ALLOWED_METHOD       = "許可されていないメソッドです。"
	HTTP_ERR_TOO_LARGE_REQUEST_BODY   = "リクエストボディのサイズが大きすぎます。"
	HTTP_ERR_FAILED_READ_REQUEST_BODY = "リクエストボディの読み取りに失敗しました。"
	HTTP_ERR_TOO_MANY_REQUESTS        = "リクエストが多すぎます。しばらく待ってから再度お試しください。"
)

// 入力関連のエラーメッセージ
const (
	INPUT_ERR_INVALID_INPUT     = "入力が不正です。"
	INPUT_ERR_INVALID_ID        = "IDが不正です。"
	INPUT_ERR_FAILED_GET_ID     = "IDの取得に失敗しました。"
	INPUT_ERR_REQUIRED_TITLE    = "タイトルは必須です。"
	INPUT_ERR_OVER_LENGTH_TITLE = "タイトルは255文字以内で入力してください。"
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
