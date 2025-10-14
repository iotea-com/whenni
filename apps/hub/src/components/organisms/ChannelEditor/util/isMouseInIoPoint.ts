import { IoPoint } from '../types/IoPoint'
import { isPointInPolygon } from './isPointInPolygon'

export const isMouseInIoPoint = (
  mouseX: number,
  mouseY: number,
  ioPoint: IoPoint,
  canvasScale: number,
): boolean => {
  const topLeftCorner = {
    x: ioPoint.coordinates.x - 5 * canvasScale,
    y: ioPoint.coordinates.y + 5 * canvasScale,
  }
  const topRightCorner = {
    x: ioPoint.coordinates.x + 5 * canvasScale,
    y: ioPoint.coordinates.y + 5 * canvasScale,
  }
  const bottomLeftCorner = {
    x: ioPoint.coordinates.x - 5 * canvasScale,
    y: ioPoint.coordinates.y - 5 * canvasScale,
  }
  const bottomRightCorner = {
    x: ioPoint.coordinates.x + 5 * canvasScale,
    y: ioPoint.coordinates.y - 5 * canvasScale,
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
