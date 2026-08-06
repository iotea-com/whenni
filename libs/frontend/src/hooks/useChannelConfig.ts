import { useMemo } from 'react'
import type { ChannelNode } from '@gruent/libs/engine/nodes/v1'

import type {
  ChannelConfig,
  ChannelEdge,
  ChannelNote,
  ChannelRuntime,
} from '@gruent/libs/engine/channels/index'

const useChannelConfig = (
  id: string,
  name: string,
  runtime: ChannelRuntime,
  nodes: Map<string, ChannelNode>,
  edges: Map<string, ChannelEdge>,
  notes: Map<string, ChannelNote>,
) => {
  const channelConfig = useMemo<ChannelConfig>(() => {
    const nodesArray = Array.from(nodes.values())
    const edgesArray = Array.from(edges.values())
    const notesArray = Array.from(notes.values())

    return {
      id,
      name,
      runtime,
      nodes: nodesArray,
      edges: edgesArray,
      notes: notesArray,
    }
  }, [id, name, runtime, nodes, edges, notes])

  return channelConfig
}

export default useChannelConfig
