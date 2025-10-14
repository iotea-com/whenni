import { ChannelEdge } from '@iotea/libs/engine/channels/channels'
import theme from '@iotea/libs/frontend/themes/tailwind'
import { IoPoint } from '../types/IoPoint'
import { ChannelNode } from '@iotea/libs/engine/nodes/v1'
import { Theme } from '@iotea/hub/stores/settingsStore'

type DrawEdgesParams = {
  context: CanvasRenderingContext2D
  width: number
  height: number
  zoom: number
  editor: HTMLCanvasElement
  edges: Map<string, ChannelEdge>
  nodes: Map<string, ChannelNode>
  ioPoints: Map<string, IoPoint>
  currentNode: ChannelNode | null
  currentEdge: ChannelEdge | null
  canvasScale: number
  dashOffset: number
  userTheme: Theme
}

const drawEdges = ({
  context,
  canvasScale,
  nodes,
  edges,
  ioPoints,
  currentNode,
  currentEdge,
  dashOffset = 0,
  userTheme,
}: DrawEdgesParams) => {
  const edgeEntries = Array.from(edges.entries())
  for (const edgeEntry of edgeEntries) {
    const edge = edgeEntry[1]
    const fromIoPoint = ioPoints.get(`${edge.from.nodeId}_${edge.from.ioId}`)
    const toIoPoint = ioPoints.get(`${edge.to.nodeId}_${edge.to.ioId}`)
    const fromNode = nodes.has(edge.from.nodeId)
    const toNode = nodes.has(edge.to.nodeId)

    // catch missing fromNode in nodes map
    if (!fromNode || !toNode || !fromIoPoint || !toIoPoint) continue

    const angle = Math.atan2(
      toIoPoint.coordinates.y - fromIoPoint.coordinates.y,
      toIoPoint.coordinates.x - fromIoPoint.coordinates.x,
    )
    const arrowLength = 10 * canvasScale
    const arrowOffset = 5 * canvasScale // Increased offset to account for arrow size

    // Calculate the point where the arrow should be (offset from the IO point)
    const arrowTip = {
      x: toIoPoint.coordinates.x,
      y: toIoPoint.coordinates.y - arrowOffset * Math.sin(angle),
    }

    // Draw the main line with bezier curve
    const startPoint = {
      x: fromIoPoint.coordinates.x,
      y: fromIoPoint.coordinates.y,
    }
    const endPoint = {
      x: arrowTip.x,
      y: arrowTip.y,
    }

    // Calculate control points for the bezier curve
    const midY = (startPoint.y + endPoint.y) / 2
    const controlPoint1 = {
      x: startPoint.x,
      y: midY,
    }
    const controlPoint2 = {
      x: endPoint.x,
      y: midY,
    }

    // Draw the curved line
    context.beginPath()
    context.moveTo(startPoint.x, startPoint.y)
    context.bezierCurveTo(
      controlPoint1.x,
      controlPoint1.y,
      controlPoint2.x,
      controlPoint2.y,
      endPoint.x,
      endPoint.y - arrowLength + canvasScale + 1,
    )

    // Determine edge color based on source node type
    const fn = nodes.get(edge.from.nodeId)

    const edgeColor = (() => {
      if (
        (currentNode && edge.from.nodeId === currentNode.id) ||
        (currentEdge && currentEdge.id === edge.id)
      ) {
        switch (fn?.metadata.type) {
          case 'source':
            return theme.colors.primary
          case 'processing':
            return theme.colors.blue[500]
          case 'action':
            return theme.colors.red[500]
          case 'conditional':
            return theme.colors.orange[500]
          case 'custom':
            return theme.colors.purple[500]
          default:
            return theme.colors.gray[500]
        }
      }

      return userTheme === 'light' ? theme.colors.gray[400] : theme.colors.gray[100]
    })()

    context.strokeStyle = edgeColor
    context.lineWidth = 1 * canvasScale

    // Apply dashed line if this edge starts from the selected node or is the selected edge
    if (
      (currentNode && edge.from.nodeId === currentNode.id) ||
      (currentEdge && currentEdge.id === edge.id)
    ) {
      context.setLineDash([2 * canvasScale, 2 * canvasScale]) // Increased dash length
      context.lineDashOffset = -(dashOffset * 0.2) * canvasScale // Slowed down the animation
    } else {
      context.setLineDash([])
    }

    context.stroke()
    context.setLineDash([]) // Reset dash pattern for other drawings

    // Draw the arrow (using the angle of the curve at the endpoint)
    const curveEndAngle = Math.atan2(endPoint.y - controlPoint2.y, endPoint.x - controlPoint2.x)

    context.beginPath()
    context.moveTo(
      endPoint.x - arrowLength * Math.cos(curveEndAngle - Math.PI / 6),
      endPoint.y - arrowLength * Math.sin(curveEndAngle - Math.PI / 6),
    )
    context.lineTo(endPoint.x, endPoint.y)
    context.lineTo(
      endPoint.x - arrowLength * Math.cos(curveEndAngle + Math.PI / 6),
      endPoint.y - arrowLength * Math.sin(curveEndAngle + Math.PI / 6),
    )
    context.fillStyle = edgeColor
    context.fill()
  }
}

export default drawEdges
