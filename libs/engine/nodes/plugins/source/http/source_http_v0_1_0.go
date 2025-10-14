// Code generated from JSON Schema using quicktype. DO NOT EDIT.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    sourceHTTPV010, err := UnmarshalSourceHTTPV010(bytes)
//    bytes, err = sourceHTTPV010.Marshal()

package codegen

import "encoding/json"

func UnmarshalSourceHTTPV010(data []byte) (SourceHTTPV01_0, error) {
	var r SourceHTTPV01_0
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SourceHTTPV01_0) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type SourceHTTPV01_0 struct {
	Method                                         *Method  `json:"method,omitempty"`
	// The timeout for the response in milliseconds         
	ResponseTimeout                                *float64 `json:"responseTimeout,omitempty"`
}

type Method string

const (
	Delete Method = "DELETE"
	Get    Method = "GET"
	Patch  Method = "PATCH"
	Post   Method = "POST"
	Put    Method = "PUT"
)
