package errors

// New 创建错误
func New(code, msg string) *AppError {
	return &AppError{CodeVal: code, MsgVal: msg}
}

// NewCategorized 创建带分类的错误
func NewCategorized(cat, code, msg string) *AppError {
	return &AppError{CodeVal: code, MsgVal: msg, CatVal: cat}
}

// Wrap 包装错误
func Wrap(err error, code, msg string) *AppError {
	if err == nil {
		return nil
	}
	return &AppError{CodeVal: code, MsgVal: msg, CauseVal: err}
}