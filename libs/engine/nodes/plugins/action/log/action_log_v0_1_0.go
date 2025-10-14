// Code generated from JSON Schema using quicktype. DO NOT EDIT.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    actionLogV010, err := UnmarshalActionLogV010(bytes)
//    bytes, err = actionLogV010.Marshal()

package codegen

import "encoding/json"

func UnmarshalActionLogV010(data []byte) (ActionLogV01_0, error) {
	var r ActionLogV01_0
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ActionLogV01_0) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ActionLogV01_0 struct {
	// The message to log                                                        
	Message                                                              string  `json:"message"`
	// The ID of the model to use if the message includes dynamic content        
	TemplateModel                                                        *string `json:"templateModel,omitempty"`
}
