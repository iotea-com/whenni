import { ChannelNodeIO } from '@gruent/libs/engine/nodes/v1'

export type IoPoint = {
  io: 'input' | 'output'
  nodeId: string
  coordinates: {
    x: number
    y: number
  }
  channelNodeIo: ChannelNodeIO
}
