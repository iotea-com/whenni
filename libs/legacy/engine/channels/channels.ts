import { ChannelNode } from "../nodes/v1/node"

export type ChannelEdge = {
  id: string
  from: ChannelEdgeConnection
  to: ChannelEdgeConnection
}

export type ChannelNote = {
  id: string
  coordinates: {
    x: number
    y: number
  }
  text: string
}

export type ChannelEdgeConnection = {
  nodeId: string
  ioId: string
}

export type ChannelRuntime = {
  size: 'small' | 'medium' | 'large'
}

export type ChannelConfig = {
  id: string
  name: string
  runtime: ChannelRuntime
  nodes: ChannelNode[]
  edges: ChannelEdge[]
  notes: ChannelNote[]
}