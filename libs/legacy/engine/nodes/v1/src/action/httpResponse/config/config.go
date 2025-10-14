package httpResponseActionNodeConfig

type HttpResponseActionNodeConfig struct {
	ResponseCode int               `json:"responseCode" validate:"required,min=100,max=599"`
	Headers      map[string]string `json:"headers"`
}
