package main

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *MessageQueueActionNode) Exec(params node.ExecParams) node.Error {
	if n.subnode == nil {
		return node.Error{
			Type:   node.FatalError,
			Reason: "a subnode was not properly routed to",
		}
	}

	return n.subnode.Exec(params)
}
