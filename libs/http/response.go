package ioteahttp

import "encoding/json"

// Structured format for every API response. Contains data, errors, and pagination.
type IoteaApiResponse struct {
	Data           any   `json:"data"`
	Errors         []any `json:"errors"`
	Page           *int  `json:"page,omitempty"`
	TotalPages     *int  `json:"totalPages,omitempty"`
	TotalResults   *int  `json:"totalResults,omitempty"`
	ResultsPerPage *int  `json:"resultsPerPage,omitempty"`
}

// Creates a new API response object for a GET (Object) route. The data parameter should be a
// unique instance of a struct.
func NewGetResponse(data any) IoteaApiResponse {
	r := IoteaApiResponse{
		Data:   data,
		Errors: nil,
	}

	return r
}

// Creates a new API response object for a GET (List) route. The data parameter should be a slice.
func NewListResponse(data any, page *int, totalPages *int, totalResults *int, resultsPerPage *int) IoteaApiResponse {
	r := IoteaApiResponse{
		Data:           data,
		Errors:         nil,
		Page:           page,
		TotalPages:     totalPages,
		TotalResults:   totalResults,
		ResultsPerPage: resultsPerPage,
	}

	return r
}

// Creates a new API response object for a POST route.
func NewCreateResponse(data any) IoteaApiResponse {
	r := IoteaApiResponse{
		Data:   data,
		Errors: nil,
	}

	return r
}

// Creates a new API response object for a PUT or PATCH route.
func NewUpdateResponse(data any, errors []any) IoteaApiResponse {
	r := IoteaApiResponse{
		Data:   data,
		Errors: errors,
	}

	return r
}

// Creates a new API response object for a DELETE route.
func NewDeleteResponse() IoteaApiResponse {
	r := IoteaApiResponse{
		Data:   nil,
		Errors: nil,
	}

	return r
}

// Marshals the API response instance to a JSON string.
func (r IoteaApiResponse) MarshalJson() ([]byte, error) {
	var errors []any
	errors = append(errors, r.Errors...)

	apiResponse := IoteaApiResponse{
		Data:       r.Data,
		Errors:     errors,
		Page:       r.Page,
		TotalPages: r.TotalPages,
	}

	marshalledResponse, err := json.Marshal(apiResponse)
	if err != nil {
		return nil, err
	}

	return marshalledResponse, err
}
