package search

import (
	ioteahttp "github.com/iotea-com/iotea/libs/http"
)

func respond(request *ioteahttp.Request[Input], output *Output) {
	request.Span.AddEvent("respond")

	page := 1
	totalPages := 1
	totalResults := len(output.Users)
	resultsPerPage := 9999
	response := ioteahttp.NewListResponse(output.Users, &page, &totalPages, &totalResults, &resultsPerPage)
	request.FiberContext.JSON(response)
}
