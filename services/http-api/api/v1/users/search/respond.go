package search

import (
	gruenthttp "github.com/ongruent/gruent/libs/http"
)

func respond(request *gruenthttp.Request[Input], output *Output) {
	request.Span.AddEvent("respond")

	page := 1
	totalPages := 1
	totalResults := len(output.Users)
	resultsPerPage := 9999
	response := gruenthttp.NewListResponse(output.Users, &page, &totalPages, &totalResults, &resultsPerPage)
	request.FiberContext.JSON(response)
}
