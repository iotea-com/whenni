package magicLinkVerify

import (
	ioteahttp "github.com/iotea-com/iotea/libs/http"
)

func contextValidate(request *ioteahttp.Request[Input]) error {
	request.Span.AddEvent("contextValidate")

	return nil
}
