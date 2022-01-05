package response

var (
	StatusSuccess = "success"
	StatusFailed  = "failed"
	ModuleName    = "bridgx/containers-cloud"
	Version       = ""
)

type ResponseBase struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func MkFailedResponse(message string) *ResponseBase {
	return &ResponseBase{
		Status:  StatusFailed,
		Message: message,
	}
}

func MkSuccessResponse(data interface{}) *ResponseBase {
	return &ResponseBase{
		Status: StatusSuccess,
		Data:   data,
	}
}

const (
	ParamError      = "参数错误"
	MetricNameError = "指标名称错误"
)
