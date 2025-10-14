package models_test

import (
	"testing"

	"github.com/iotea-com/iotea/libs/engine/dependencies/models"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	t.Run("valid model with primitive attributes", func(t *testing.T) {
		model := models.Model{
			Attributes: map[string]models.Attribute{
				"attribute1": {
					Id:   "attribute1",
					Key:  "myString",
					Type: "string",
				},
				"attribute2": {
					Id:   "attribute2",
					Key:  "myNumber",
					Type: "number",
				},
			},
		}

		err := model.Validate()
		require.NoError(t, err)
	})

	t.Run("valid nested object", func(t *testing.T) {
		model := models.Model{
			Attributes: map[string]models.Attribute{
				"parent": {
					Id:   "parent",
					Key:  "parent",
					Type: "object",
				},
				"child": {
					Id:       "child",
					Key:      "child",
					Type:     "string",
					ParentId: "parent",
				},
			},
		}

		err := model.Validate()
		require.NoError(t, err)
	})

	t.Run("valid array with primitive subtype", func(t *testing.T) {
		model := models.Model{
			Attributes: map[string]models.Attribute{
				"array": {
					Id:      "array",
					Key:     "numbers",
					Type:    "array",
					Subtype: "number",
				},
			},
		}

		err := model.Validate()
		require.NoError(t, err)
	})

	t.Run("valid array with object reference", func(t *testing.T) {
		model := models.Model{
			Attributes: map[string]models.Attribute{
				"array": {
					Id:      "array",
					Key:     "items",
					Type:    "array",
					Subtype: "item",
				},
				"item": {
					Id:   "item",
					Key:  "item",
					Type: "object",
				},
				"itemAttribute": {
					Id:       "itemAttribute",
					Key:      "name",
					Type:     "string",
					ParentId: "item",
				},
			},
		}

		err := model.Validate()
		require.NoError(t, err)
	})

	t.Run("invalid attribute key", func(t *testing.T) {
		model := models.Model{
			Attributes: map[string]models.Attribute{
				"attribute1": {
					Id:   "attribute1",
					Key:  "123invalid",
					Type: "string",
				},
			},
		}

		err := model.Validate()
		expectedErr := "invalid attribute '123invalid' (attribute1) in root model: Key: 'Attribute.Key' Error:Attribute validation for 'Key' failed on the 'isValidAttributeKey' tag"
		require.EqualError(t, err, expectedErr)
	})

	t.Run("duplicate keys in same level", func(t *testing.T) {
		model := models.Model{
			Attributes: map[string]models.Attribute{
				"parent": {
					Id:   "parent",
					Key:  "parent",
					Type: "object",
				},
				"child1": {
					Id:       "child1",
					Key:      "duplicate",
					Type:     "string",
					ParentId: "parent",
				},
				"child2": {
					Id:       "child2",
					Key:      "duplicate",
					Type:     "string",
					ParentId: "parent",
				},
			},
		}

		err := model.Validate()
		require.Error(t, err)
		require.Contains(t, err.Error(), "duplicate key")
	})

	t.Run("empty object", func(t *testing.T) {
		model := models.Model{
			Attributes: map[string]models.Attribute{
				"empty": {
					Id:   "empty",
					Key:  "empty",
					Type: "object",
				},
			},
		}

		err := model.Validate()
		expectedErr := "attribute 'empty' (empty): objects must have at least one child attribute"
		require.EqualError(t, err, expectedErr)
	})

	t.Run("array without subtype", func(t *testing.T) {
		model := models.Model{
			Attributes: map[string]models.Attribute{
				"array": {
					Id:   "array",
					Key:  "array",
					Type: "array",
				},
			},
		}

		err := model.Validate()
		expectedErr := "invalid attribute 'array' (array) in root model: Key: 'Attribute.Subtype' Error:Attribute validation for 'Subtype' failed on the 'required_if' tag"
		require.EqualError(t, err, expectedErr)
	})

	t.Run("invalid array subtype reference", func(t *testing.T) {
		model := models.Model{
			Attributes: map[string]models.Attribute{
				"array": {
					Id:      "array",
					Key:     "array",
					Type:    "array",
					Subtype: "nonexistent",
				},
			},
		}

		err := model.Validate()
		expectedErr := "attribute 'array' (array): invalid array subtype reference 'nonexistent'"
		require.EqualError(t, err, expectedErr)
	})

	t.Run("duplicate keys in different levels", func(t *testing.T) {
		model := models.Model{
			Attributes: map[string]models.Attribute{
				"id1": {
					Id:   "id1",
					Key:  "object",
					Type: "object",
				},
				"id2": {
					Id:       "id2",
					Key:      "duplicate",
					Type:     "string",
					ParentId: "id1",
				},
				"id3": {
					Id:       "id3",
					Key:      "duplicate",
					Type:     "string",
					ParentId: "id1",
				},
			},
		}

		err := model.Validate()
		require.Error(t, err)
		require.Contains(t, err.Error(), "duplicate key")
	})

	t.Run("circular reference", func(t *testing.T) {
		model := models.Model{
			Attributes: map[string]models.Attribute{
				"obj1": {
					Id:       "obj1",
					Key:      "obj1",
					Type:     "object",
					ParentId: "circularDependency1",
				},
				"circularDependency1": {
					Id:       "circularDependency1",
					Key:      "circular",
					Type:     "object",
					ParentId: "obj1",
				},
				"string1": {
					Id:       "string1",
					Key:      "string1",
					Type:     "string",
					ParentId: "obj1",
				},
				"string2": {
					Id:       "string2",
					Key:      "string2",
					Type:     "string",
					ParentId: "circularDependency1",
				},
			},
		}

		err := model.Validate()
		require.ErrorContains(t, err, "circular reference detected")
	})
}
