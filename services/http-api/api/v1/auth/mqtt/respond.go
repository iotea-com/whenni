package mqtt

import (
	ioteahttp "github.com/iotea-com/iotea/libs/http"
)

type ResponseBody struct {
	Result string   `json:"result"`
	Errors []string `json:"errors"`
}

func respond(request *ioteahttp.Request[Input], output *Output) {
	request.Span.AddEvent("respond")

	request.FiberContext.Status(200).JSON(ResponseBody{
		Result: output.Result,
		Errors: nil,
	})
}
