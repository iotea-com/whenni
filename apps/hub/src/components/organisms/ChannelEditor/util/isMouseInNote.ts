import { ChannelNote } from '@iotea/libs/engine/channels/index'
import { calculateNoteDisplay } from './calculateNoteDisplay'
import { isPointInPolygon } from './isPointInPolygon'

export const isMouseInNote = (
  mouseX: number,
  mouseY: number,
  note: ChannelNote,
  canvasScale: number,
  pan: { x: number; y: number },
): boolean => {
  const { x, y, width, height } = calculateNoteDisplay(note, canvasScale, pan)

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
