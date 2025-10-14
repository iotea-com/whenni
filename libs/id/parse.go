package id

import (
	"fmt"
)

type IdType string

const (
	IdTypeUser             IdType = "u"
	IdTypeModel            IdType = "m"
	IdTypeThing            IdType = "t"
	IdTypeChannel          IdType = "c"
	IdTypeApiKey           IdType = "tea_"
	IdTypePermissionSet    IdType = "p"
	IdTypeOrganization     IdType = "o"
	IdTypeSpace            IdType = "s"
	IdTypeChannelExecution IdType = "x"
	IdTypeNode             IdType = "n"
)

func Parse(id string) (IdType, error) {
	if len(id) < 2 {
		return "", fmt.Errorf("invalid ID: %s", id)
	}

	if id[0] == 't' {
		if id[:5] == string(IdTypeApiKey) {
			return IdTypeApiKey, nil
		}
	}

	switch string(id[0]) {
	case string(IdTypeOrganization):
		return IdTypeOrganization, nil
	case string(IdTypeSpace):
		return IdTypeSpace, nil
	case string(IdTypeUser):
		return IdTypeUser, nil
	case string(IdTypeModel):
		return IdTypeModel, nil
	case string(IdTypeThing):
		return IdTypeThing, nil
	case string(IdTypeChannel):
		return IdTypeChannel, nil
	case string(IdTypeNode):
		return IdTypeNode, nil
	case string(IdTypePermissionSet):
		return IdTypePermissionSet, nil
	case string(IdTypeChannelExecution):
		return IdTypeChannelExecution, nil
	}

	return "", fmt.Errorf("invalid ID: %s", id)
}
