// Code generated from JSON Schema using quicktype. DO NOT EDIT.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    conditionalThresholdV010, err := UnmarshalConditionalThresholdV010(bytes)
//    bytes, err = conditionalThresholdV010.Marshal()

package codegen

import "encoding/json"

func UnmarshalConditionalThresholdV010(data []byte) (ConditionalThresholdV01_0, error) {
	var r ConditionalThresholdV01_0
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ConditionalThresholdV01_0) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ConditionalThresholdV01_0 struct {
	// The threshold conditions to evaluate                                          
	Conditions                                                       []ConfigSchema  `json:"conditions,omitempty"`
	// The logical operator to use for evaluating multiple conditions                
	LogicalOperator                                                  LogicalOperator `json:"logicalOperator"`
}

type ConfigSchema struct {
	// The ID of the attribute to evaluate                       
	AttributeID                                        *string   `json:"attributeId,omitempty"`
	// The operator to use for evaluating the condition          
	Operator                                           *Operator `json:"operator,omitempty"`
	// The value to compare against                              
	Value                                              *float64  `json:"value,omitempty"`
}

// The operator to use for evaluating the condition
type Operator string

const (
	EqualTo            Operator = "equal-to"
	GreaterThan        Operator = "greater-than"
	GreaterThanEqualTo Operator = "greater-than-equal-to"
	LessThan           Operator = "less-than"
	LessThanEqualTo    Operator = "less-than-equal-to"
	NotEqualTo         Operator = "not-equal-to"
)

// The logical operator to use for evaluating multiple conditions
type LogicalOperator string

const (
	And LogicalOperator = "AND"
	Or  LogicalOperator = "OR"
)
