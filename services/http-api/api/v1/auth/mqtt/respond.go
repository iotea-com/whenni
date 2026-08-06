package mqtt

import (
	gruenthttp "github.com/ongruent/gruent/libs/http"
)

type ResponseBody struct {
	Result string   `json:"result"`
	Errors []string `json:"errors"`
}

func respond(request *gruenthttp.Request[Input], output *Output) {
	request.Span.AddEvent("respond")

	request.FiberContext.Status(200).JSON(ResponseBody{
		Result: output.Result,
		Errors: nil,
	})
}
