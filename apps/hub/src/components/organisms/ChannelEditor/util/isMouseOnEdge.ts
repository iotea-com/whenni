import { ChannelEdge } from '@gruent/libs/engine/channels/index'
import { IoPoint } from '../types/IoPoint'
import { isPointInPolygon } from './isPointInPolygon'

export const isMouseOnEdge = (
  mouseX: number,
  mouseY: number,
  edge: ChannelEdge,
  ioPoints: IoPoint[],
  canvasScale: number,
): ChannelEdge | null => {
  // get IO points for edge
  const fromNodeIoPoint = ioPoints.filter((ioPoint) => {
    // check if both the nodeId and the ioId match up
    if (ioPoint.nodeId === edge.from.nodeId && ioPoint.channelNodeIo.id === edge.from.ioId)
      return true

    return false
  })[0]

  const toNodeIoPoint = ioPoints.filter((ioPoint) => {
    // check if both the nodeId and the ioId match up
    if (ioPoint.nodeId === edge.to.nodeId && ioPoint.channelNodeIo.id === edge.to.ioId) return true

    return false
  })[0]

  if (!fromNodeIoPoint || !toNodeIoPoint) return null

  // calculate slope of the line
  const deltaY = toNodeIoPoint.coordinates.y - fromNodeIoPoint.coordinates.y
  const deltaX = toNodeIoPoint.coordinates.x - fromNodeIoPoint.coordinates.x
  const slope = deltaY / deltaX

  // calculate perpendicular slope
  const perpSlope = -1 / slope

  // function to calculate y using the line equation y = mx + c
  const getY = (x, slope, x1, y1) => slope * (x - x1) + y1

  // calculate the corners of the hitbox
  const hitboxWidth = 10 // margin size for easier clicking
  const offset = hitboxWidth / Math.sqrt(1 + Math.pow(perpSlope, 2))
  const topRightCorner = {
    x: toNodeIoPoint.coordinates.x - offset,
    y: getY(
      toNodeIoPoint.coordinates.x - offset,
      perpSlope,
      toNodeIoPoint.coordinates.x,
      toNodeIoPoint.coordinates.y,
    ),
  }

  const topLeftCorner = {
    x: toNodeIoPoint.coordinates.x + offset,
    y: getY(
      toNodeIoPoint.coordinates.x + offset,
      perpSlope,
      toNodeIoPoint.coordinates.x,
      toNodeIoPoint.coordinates.y,
    ),
  }

  const bottomLeftCorner = {
    x: fromNodeIoPoint.coordinates.x + offset,
    y: getY(
      fromNodeIoPoint.coordinates.x + offset,
      perpSlope,
      fromNodeIoPoint.coordinates.x,
      fromNodeIoPoint.coordinates.y,
    ),
  }

  const bottomRightCorner = {
    x: fromNodeIoPoint.coordinates.x - offset,
    y: getY(
      fromNodeIoPoint.coordinates.x - offset,
      perpSlope,
      fromNodeIoPoint.coordinates.x,
      fromNodeIoPoint.coordinates.y,
    ),
  }

  // check if point is in the polygon
  if (
    isPointInPolygon(
      { x: mouseX, y: mouseY },
      [topLeftCorner, topRightCorner, bottomRightCorner, bottomLeftCorner],
      canvasScale,
    )
  )
    return edge

  return null
}
