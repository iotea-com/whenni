package gruentpermissions

import (
	"fmt"
	"strings"
)

type Namespace string
type Action string
type Permission string

var ActionAll Action = "*"
var ActionList Action = "list"
var ActionCreate Action = "create"
var ActionGet Action = "get"
var ActionUpdate Action = "update"
var ActionDelete Action = "delete"

func Marshal(namespace Namespace, action Action) string {
	return fmt.Sprintf("%s:%s", namespace, action)
}

func Unmarshal(s string) (*Namespace, *Action, error) {
	permissionParts := strings.Split(s, ":")
	if len(permissionParts) != 2 {
		err := fmt.Errorf("invalid permission string")
		return nil, nil, err
	}

	// TODO: validate values of namespace and action

	namespace := Namespace(permissionParts[0])
	action := Action(permissionParts[1])
	return &namespace, &action, nil
}
