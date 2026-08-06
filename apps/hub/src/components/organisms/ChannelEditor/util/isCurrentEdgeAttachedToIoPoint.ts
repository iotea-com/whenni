import { ChannelEdge } from '@gruent/libs/engine/channels/index'
import { IoPoint } from '../types/IoPoint'

export const isCurrentEdgeAttachedToIoPoint = (
  currentEdge: ChannelEdge | null,
  ioPoint: IoPoint,
) => {
  // always return false if there's no edge selected
  if (!currentEdge) return false

  // check if the nodeId of the edge matches the nodeId of the ioPoint
  if (currentEdge.from.nodeId !== ioPoint.nodeId && currentEdge.to.nodeId !== ioPoint.nodeId)
    return false

  // check if the ioIds (which are unique per node) match between the edge and the ioPoint
  if (currentEdge.from.ioId === ioPoint.channelNodeIo.id) return true
  if (currentEdge.to.ioId === ioPoint.channelNodeIo.id) return true

  // default return false
  return false
}
