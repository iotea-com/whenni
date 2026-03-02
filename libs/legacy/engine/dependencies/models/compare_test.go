package models_test

import (
	"encoding/json"
	"testing"

	"github.com/iotea-com/iotea/libs/legacy/engine/dependencies/models"
	"github.com/stretchr/testify/require"
)

func TestCompare(t *testing.T) {
	t.Run("valid simple schema", func(t *testing.T) {
		schema := models.Model{
			Attributes: map[string]models.Attribute{
				"attribute1": {
					Id:       "attribute1",
					Key:      "myString",
					Type:     "string",
					Required: true,
				},
				"attribute2": {
					Id:       "attribute2",
					Key:      "myNumber",
					Type:     "number",
					Required: true,
				},
			},
		}
		data := map[string]json.RawMessage{
			"myString": json.RawMessage(`"test string"`),
			"myNumber": json.RawMessage(`42`),
		}

		err := schema.Validate()
		if err != nil {
			t.Errorf("Schema.Validate() error = %s", err)
		}

		err = schema.Compare(data)
		require.NoError(t, err)
	})

	t.Run("invalid json", func(t *testing.T) {
		schema := models.Model{
			Attributes: map[string]models.Attribute{
				"attribute1": {
					Id:   "attribute1",
					Key:  "myString",
					Type: "string",
				},
			},
		}
		data := map[string]json.RawMessage{
			"myString": json.RawMessage(`"test str`),
		}

		err := schema.Validate()
		if err != nil {
			t.Errorf("Schema.Validate() error = %s", err)
		}

		err = schema.Compare(data)
		expectedErr := "error unmarshaling field 'myString' (attribute1): unexpected end of JSON input"
		require.EqualError(t, err, expectedErr)
	})

	t.Run("missing required field", func(t *testing.T) {
		schema := models.Model{
			Attributes: map[string]models.Attribute{
				"attribute1": {
					Id:       "attribute1",
					Key:      "myString",
					Type:     "string",
					Required: true,
				},
			},
		}
		data := map[string]json.RawMessage{}

		err := schema.Validate()
		if err != nil {
			t.Errorf("Schema.Validate() error = %s", err)
		}

		err = schema.Compare(data)
		expectedErr := "expected field 'myString' (attribute1) is missing in input data"
		require.EqualError(t, err, expectedErr)
	})

	t.Run("invalid string type", func(t *testing.T) {
		schema := models.Model{
			Attributes: map[string]models.Attribute{
				"attribute1": {
					Id:       "attribute1",
					Key:      "myString",
					Type:     "string",
					Required: true,
				},
			},
		}
		data := map[string]json.RawMessage{
			"myString": json.RawMessage(`42`),
		}

		err := schema.Validate()
		if err != nil {
			t.Errorf("Schema.Validate() error = %s", err)
		}

		err = schema.Compare(data)
		expectedErr := "attribute 'myString' (attribute1) is not a string - got 42"
		require.EqualError(t, err, expectedErr)
	})

	t.Run("invalid number type", func(t *testing.T) {
		schema := models.Model{
			Attributes: map[string]models.Attribute{
				"attribute1": {
					Id:       "attribute1",
					Key:      "myNumber",
					Type:     "number",
					Required: true,
				},
			},
		}
		data := map[string]json.RawMessage{
			"myNumber": json.RawMessage(`"test string"`),
		}

		err := schema.Validate()
		if err != nil {
			t.Errorf("Schema.Validate() error = %s", err)
		}

		err = schema.Compare(data)
		expectedErr := "attribute 'myNumber' (attribute1) is not a number - got test string"
		require.EqualError(t, err, expectedErr)
	})

	t.Run("invalid boolean type", func(t *testing.T) {
		schema := models.Model{
			Attributes: map[string]models.Attribute{
				"attribute1": {
					Id:       "attribute1",
					Key:      "myBool",
					Type:     "boolean",
					Required: true,
				},
			},
		}
		data := map[string]json.RawMessage{
			"myBool": json.RawMessage(`"test string"`),
		}

		err := schema.Validate()
		if err != nil {
			t.Errorf("Schema.Validate() error = %s", err)
		}

		err = schema.Compare(data)
		expectedErr := "attribute 'myBool' (attribute1) is not a boolean - got test string"
		require.EqualError(t, err, expectedErr)
	})

	t.Run("invalid object type - wrong subfield type", func(t *testing.T) {
		schema := models.Model{
			Attributes: map[string]models.Attribute{
				"object1": {
					Id:       "object1",
					Key:      "myObject",
					Type:     "object",
					Required: true,
				},
				"subfield1": {
					Id:       "subfield1",
					Key:      "nestedString",
					Type:     "string",
					ParentId: "object1",
					Required: true,
				},
			},
		}
		data := map[string]json.RawMessage{
			"myObject": json.RawMessage(`{"nestedString": 1}`),
		}

		err := schema.Validate()
		if err != nil {
			t.Errorf("Schema.Validate() error = %s", err)
		}

		err = schema.Compare(data)
		expectedErr := "validation failed for object myObject (object1): attribute 'nestedString' (subfield1) is not a string - got 1"
		require.EqualError(t, err, expectedErr)
	})

	t.Run("invalid array - wrong subtype", func(t *testing.T) {
		schema := models.Model{
			Attributes: map[string]models.Attribute{
				"attribute1": {
					Id:       "attribute1",
					Key:      "myArray",
					Type:     "array",
					Subtype:  "number",
					Required: true,
				},
			},
		}
		data := map[string]json.RawMessage{
			"myArray": json.RawMessage(`"test string"`),
		}

		err := schema.Validate()
		if err != nil {
			t.Errorf("Schema.Validate() error = %s", err)
		}

		err = schema.Compare(data)
		expectedErr := "attribute 'myArray' (attribute1) is not a valid number array: json: cannot unmarshal string into Go value of type []float64"
		require.EqualError(t, err, expectedErr)
	})

	t.Run("invalid object array - wrong type", func(t *testing.T) {
		schema := models.Model{
			Attributes: map[string]models.Attribute{
				"attribute1": {
					Id:      "attribute1",
					Key:     "myArray",
					Type:    "array",
					Subtype: "object",
					ArrayObjectAttributes: map[string]models.Attribute{
						"subfield1": {
							Id:       "subfield1",
							Key:      "nestedString",
							Type:     "string",
							Required: true,
						},
					},
					Required: true,
				},
			},
		}
		data := map[string]json.RawMessage{
			"myArray": json.RawMessage(`[{"nestedString": 1}]`),
		}

		err := schema.Validate()
		if err != nil {
			t.Errorf("Schema.Validate() error = %s", err)
		}

		err = schema.Compare(data)
		expectedErr := "validation failed for object 0 in array myArray (attribute1): attribute 'nestedString' (subfield1) is not a string - got 1"
		require.EqualError(t, err, expectedErr)
	})

	t.Run("valid nested object", func(t *testing.T) {
		schema := models.Model{
			Attributes: map[string]models.Attribute{
				"attribute1": {
					Id:   "attribute1",
					Key:  "attribute1",
					Type: "string",
				},
				"obj1": {
					Id:       "obj1",
					Key:      "myObject",
					Type:     "object",
					Required: true,
				},
				"subfield1": {
					Id:       "subfield1",
					Key:      "nestedString",
					Type:     "string",
					ParentId: "obj1",
					Required: true,
				},
			},
		}
		data := map[string]json.RawMessage{
			"attribute1": json.RawMessage(`"nested value"`),
			"myObject": json.RawMessage(`{
				"nestedString": "nested value"
			}`),
		}

		err := schema.Validate()
		if err != nil {
			t.Errorf("Schema.Validate() error = %s", err)
		}

		err = schema.Compare(data)
		require.NoError(t, err)
	})

	t.Run("valid array of strings", func(t *testing.T) {
		schema := models.Model{
			Attributes: map[string]models.Attribute{
				"arr1": {
					Id:       "arr1",
					Key:      "myArray",
					Type:     "array",
					Subtype:  "string",
					Required: true,
				},
			},
		}
		data := map[string]json.RawMessage{
			"myArray": json.RawMessage(`["value1", "value2"]`),
		}

		err := schema.Validate()
		if err != nil {
			t.Errorf("Schema.Validate() error = %s", err)
		}

		err = schema.Compare(data)
		require.NoError(t, err)
	})

	t.Run("valid array of numbers", func(t *testing.T) {
		schema := models.Model{
			Attributes: map[string]models.Attribute{
				"arr1": {
					Id:       "arr1",
					Key:      "myArray",
					Type:     "array",
					Subtype:  "number",
					Required: true,
				},
			},
		}
		data := map[string]json.RawMessage{
			"myArray": json.RawMessage(`[1, 2, 3.14]`),
		}

		err := schema.Validate()
		if err != nil {
			t.Errorf("Schema.Validate() error = %s", err)
		}

		err = schema.Compare(data)
		require.NoError(t, err)
	})

	t.Run("valid array of booleans", func(t *testing.T) {
		schema := models.Model{
			Attributes: map[string]models.Attribute{
				"arr1": {
					Id:       "arr1",
					Key:      "myArray",
					Type:     "array",
					Subtype:  "boolean",
					Required: true,
				},
			},
		}
		data := map[string]json.RawMessage{
			"myArray": json.RawMessage(`[true, false, true]`),
		}

		err := schema.Validate()
		if err != nil {
			t.Errorf("Schema.Validate() error = %s", err)
		}

		err = schema.Compare(data)
		require.NoError(t, err)
	})

	t.Run("valid array of objects", func(t *testing.T) {
		schema := models.Model{
			Attributes: map[string]models.Attribute{
				"arr1": {
					Id:       "arr1",
					Key:      "myArray",
					Type:     "array",
					Subtype:  "object1",
					Required: true,
				},
				"object1": {
					Id:   "object1",
					Key:  "myObject",
					Type: "object",
				},
				"subfield1": {
					Id:       "subfield1",
					Key:      "nestedString",
					Type:     "string",
					ParentId: "object1",
					Required: true,
				},
			},
		}
		data := map[string]json.RawMessage{
			"myArray": json.RawMessage(`[{"nestedString": "value1"}, {"nestedString": "value2"}]`),
		}
		err := schema.Validate()
		if err != nil {
			t.Errorf("Schema.Validate() error = %s", err)
		}

		err = schema.Compare(data)
		require.NoError(t, err)
	})

	t.Run("invalid array of strings", func(t *testing.T) {
		schema := models.Model{
			Attributes: map[string]models.Attribute{
				"arr1": {
					Id:       "arr1",
					Key:      "myArray",
					Type:     "array",
					Subtype:  "string",
					Required: true,
				},
			},
		}
		data := map[string]json.RawMessage{
			"myArray": json.RawMessage(`[1, 2, 3, "string"]`),
		}

		err := schema.Validate()
		if err != nil {
			t.Errorf("Schema.Validate() error = %s", err)
		}

		err = schema.Compare(data)
		expectedErr := "attribute 'myArray' (arr1) is not a valid string array: json: cannot unmarshal number into Go value of type string"
		require.EqualError(t, err, expectedErr)
	})

	t.Run("invalid array of numbers", func(t *testing.T) {
		schema := models.Model{
			Attributes: map[string]models.Attribute{
				"arr1": {
					Id:       "arr1",
					Key:      "numbers",
					Type:     "array",
					Subtype:  "number",
					Required: true,
				},
			},
		}
		data := map[string]json.RawMessage{
			"numbers": json.RawMessage(`[1, 2, 3.14, "string"]`),
		}

		err := schema.Validate()
		if err != nil {
			t.Errorf("Schema.Validate() error = %s", err)
		}

		err = schema.Compare(data)
		expectedErr := "attribute 'numbers' (arr1) is not a valid number array: json: cannot unmarshal string into Go value of type float64"
		require.EqualError(t, err, expectedErr)
	})

	t.Run("invalid array of booleans", func(t *testing.T) {
		schema := models.Model{
			Attributes: map[string]models.Attribute{
				"arr1": {
					Id:       "arr1",
					Key:      "myArray",
					Type:     "array",
					Subtype:  "boolean",
					Required: true,
				},
			},
		}
		data := map[string]json.RawMessage{
			"myArray": json.RawMessage(`[true, false, true, "string"]`),
		}

		err := schema.Validate()
		if err != nil {
			t.Errorf("Schema.Validate() error = %s", err)
		}

		err = schema.Compare(data)
		expectedErr := "attribute 'myArray' (arr1) is not a valid boolean array: json: cannot unmarshal string into Go value of type bool"
		require.EqualError(t, err, expectedErr)
	})

	t.Run("invalid array type", func(t *testing.T) {
		schema := models.Model{
			Attributes: map[string]models.Attribute{
				"arr1": {
					Id:       "arr1",
					Key:      "myArray",
					Type:     "array",
					Subtype:  "string",
					Required: true,
				},
			},
		}
		data := map[string]json.RawMessage{
			"myArray": json.RawMessage(`[1, 2, 3]`),
		}

		err := schema.Validate()
		if err != nil {
			t.Errorf("Schema.Validate() error = %s", err)
		}

		err = schema.Compare(data)
		expectedErr := "attribute 'myArray' (arr1) is not a valid string array: json: cannot unmarshal number into Go value of type string"
		require.EqualError(t, err, expectedErr)
	})

	t.Run("optional field missing", func(t *testing.T) {
		schema := models.Model{
			Attributes: map[string]models.Attribute{
				"attribute1": {
					Id:       "attribute1",
					Key:      "optionalField",
					Type:     "string",
					Required: false,
				},
			},
		}
		data := map[string]json.RawMessage{}

		err := schema.Validate()
		if err != nil {
			t.Errorf("Schema.Validate() error = %s", err)
		}

		err = schema.Compare(data)
		require.NoError(t, err)
	})

	t.Run("valid boolean field", func(t *testing.T) {
		schema := models.Model{
			Attributes: map[string]models.Attribute{
				"attribute1": {
					Id:       "attribute1",
					Key:      "myBool",
					Type:     "boolean",
					Required: true,
				},
			},
		}
		data := map[string]json.RawMessage{
			"myBool": json.RawMessage(`true`),
		}

		err := schema.Validate()
		if err != nil {
			t.Errorf("Schema.Validate() error = %s", err)
		}

		err = schema.Compare(data)
		require.NoError(t, err)
	})
}
