package services

type Response struct {
	Code int `json:"code"`
	Data any `json:"data"`
}

const (
	CodeContentNotFound = 902000
	CodeAccessDenied    = 902001
	CodeUnauthorized    = 902002
	CodeInvalidParams   = 902003
	CodeConflict        = 902004
	CodeInternalError   = 902005
)

func resOk(data any) (int, Response) {
	return 200, Response{
		Code: 0,
		Data: data,
	}
}

func unauthorized() (int, Response) {
	return 401, Response{
		Code: CodeUnauthorized,
	}
}

func notFound(data any) (int, Response) {
	return 404, Response{
		Code: CodeContentNotFound,
		Data: data,
	}
}

func invalidParams(data any) (int, Response) {
	return 400, Response{
		Code: CodeInvalidParams,
		Data: data,
	}
}

func internalError(data any) (int, Response) {
	return 500, Response{
		Code: CodeInternalError,
		Data: data,
	}
}
