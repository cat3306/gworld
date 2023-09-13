package engine

const (
	ResponseOK  = 200
	ResponseErr = -1
)

type JsonResponse struct {
	Code int         `json:"Code"`
	Msg  string      `json:"Msg"`
	Data interface{} `json:"Data,omitempty"`
}

func JsonRspOK(data interface{}) *JsonResponse {
	return &JsonResponse{
		Code: ResponseOK,
		Data: data,
		Msg:  "",
	}
}
func JsonRspErr(msg string) *JsonResponse {
	return &JsonResponse{
		Code: ResponseOK,
		Msg:  msg,
	}
}
