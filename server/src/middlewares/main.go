package middlewares

type ContextParam struct {
	Err     error
	Message string
	Data    interface{}
}
