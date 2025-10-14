import { Theme } from '@iotea/hub/stores/settingsStore'

type DrawGridParams = {
  context: CanvasRenderingContext2D
  width: number
  height: number
  gridSize: number
  canvasScale: number
  pan: { x: number; y: number }
  userTheme: Theme
}

const drawGrid = ({
  context,
  width,
  height,
  gridSize,
  canvasScale,
  pan,
  userTheme,
}: DrawGridParams) => {
  context.fillStyle = userTheme === 'light' ? '#e8e8e8' : '#3a3a3a'
  const startX = -(pan.x % (gridSize * canvasScale))
  const startY = -(pan.y % (gridSize * canvasScale))

  for (let x = startX; x < width; x += gridSize * canvasScale) {
    for (let y = startY; y < height; y += gridSize * canvasScale) {
      context.beginPath()
      context.arc(x, y, 1 * canvasScale, 0, 2 * Math.PI) // draws a circle of radius 1
      context.fill()
    }
  }
}

export default drawGrid
