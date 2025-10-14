import { Theme } from '@iotea/hub/stores/settingsStore'
import theme from '@iotea/libs/frontend/themes/tailwind'

type DrawBlankCanvasParams = {
  context: CanvasRenderingContext2D
  width: number
  height: number
  zoom: number
  editor: HTMLCanvasElement
  userTheme: Theme
}

const drawBlankCanvas = ({
  context,
  width,
  height,
  zoom,
  editor,
  userTheme,
}: DrawBlankCanvasParams) => {
  // Apply zoom transformation
  context.save()
  context.translate(width / 2, height / 2)
  context.scale(zoom, zoom)
  context.translate(-width / 2, -height / 2)

  // Reset the canvas
  context.clearRect(0, 0, editor.width, editor.height)
  context.beginPath()
  context.rect(0, 0, editor.width, editor.height)
  context.fillStyle = userTheme === 'light' ? theme.colors.gray[50] : '#1c1d1c'
  context.fillRect(0, 0, editor.width, editor.height)

  // Add anti-aliasing
  context.imageSmoothingEnabled = true
  context.textRendering = 'optimizeLegibility'
}

export default drawBlankCanvas
