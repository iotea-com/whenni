package signin

import (
	gruenthttp "github.com/ongruent/gruent/libs/http"
)

func contextValidate(request *gruenthttp.Request[Input]) error {
	request.Span.AddEvent("contextValidate")

	return nil
}
