import { ChannelNode } from '@gruent/libs/engine/nodes/v1'
import { calculateNodeDisplay } from './calculateNodeDisplay'
import { isPointInPolygon } from './isPointInPolygon'

export const isMouseInNode = (
  mouseX: number,
  mouseY: number,
  node: ChannelNode,
  canvasScale: number,
  pan: { x: number; y: number },
): boolean => {
  const { x, y, width, height } = calculateNodeDisplay(node, canvasScale, pan)

  const margin = 2

  const topLeftCorner = {
    x: x - margin,
    y: y - margin,
  }
  const topRightCorner = {
    x: x + width + margin,
    y: y - margin,
  }
  const bottomLeftCorner = {
    x: x - margin,
    y: y + height + margin,
  }
  const bottomRightCorner = {
    x: x + width + margin,
    y: y + height + margin,
  }

  if (
    isPointInPolygon(
      { x: mouseX, y: mouseY },
      [topLeftCorner, topRightCorner, bottomRightCorner, bottomLeftCorner],
      canvasScale,
    )
  )
    return true

  return false
}
