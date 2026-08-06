import { ChannelEdge } from '@gruent/libs/engine/channels/channels'
import { ChannelNode } from '@gruent/libs/engine/nodes/v1'
import { IoPoint } from '../types/IoPoint'
import { Theme } from '@gruent/hub/stores/settingsStore'
import { calculateNodeDisplay } from '../util/calculateNodeDisplay'
import { isCurrentEdgeAttachedToIoPoint } from '../util/isCurrentEdgeAttachedToIoPoint'
import theme from '@gruent/libs/frontend/themes/tailwind'
import nodeIcons from '../util/nodeIcons'

type DrawNodesParams = {
  context: CanvasRenderingContext2D
  width: number
  height: number
  zoom: number
  editor: HTMLCanvasElement
  nodes: Map<string, ChannelNode>
  edges: Map<string, ChannelEdge>
  ioPoints: Map<string, IoPoint>
  currentNode: ChannelNode | null
  currentEdge: ChannelEdge | null
  canvasScale: number
  pan: { x: number; y: number }
  userTheme: Theme
  upsertIoPoint: (key: string, ioPoint: IoPoint) => void
  nodeErrors: Record<string, string[]> | null
}

const drawNodes = ({
  context,
  nodes,
  edges,
  ioPoints,
  currentNode,
  currentEdge,
  canvasScale,
  pan,
  userTheme,
  upsertIoPoint,
  nodeErrors,
}: DrawNodesParams) => {
  for (const nodeEntry of nodes) {
    const nodeId = nodeEntry[0]
    const node = nodeEntry[1]

    const {
      x,
      y,
      width,
      height,
      backgroundColor,
      borderColor,
      borderWidth,
      selectedBorderColor,
      selectedBorderWidth,
      accentColor,
      labelTextColor,
      typeTextColor,
      typeIcon,
      labelIcon,
    } = calculateNodeDisplay(node, canvasScale, pan, userTheme)

    // Fill rectangle background
    context.beginPath()
    context.shadowColor = 'rgba(1, 1, 3, 0.05)'
    context.shadowBlur = 8 * canvasScale
    context.shadowOffsetX = 0
    context.shadowOffsetY = 4 * canvasScale
    context.roundRect(x, y, width, height, 3)
    context.fillStyle = backgroundColor
    context.fill()

    // Reset shadow settings before drawing borders
    context.shadowColor = 'transparent'
    context.shadowBlur = 0
    context.shadowOffsetX = 0
    context.shadowOffsetY = 0

    // Apply stroke if this is the current node
    if (currentNode && nodeId === currentNode.id) {
      context.strokeStyle = selectedBorderColor
      context.lineWidth = selectedBorderWidth * canvasScale
      context.stroke()
    } else {
      context.strokeStyle = borderColor
      context.lineWidth = borderWidth * canvasScale
      context.stroke()
    }

    // Fill type tab to the left of the node
    context.beginPath()
    context.rect(x - 16.5 * canvasScale, y, 16.5 * canvasScale, height)
    if (currentNode && nodeId === currentNode.id) {
      context.fillStyle = selectedBorderColor
    } else {
      context.fillStyle = accentColor
    }
    context.stroke()
    context.fill()

    // Fill type icon
    const typeIconSize = 16
    context.save()
    context.translate(
      x - 13.5 * canvasScale,
      y + (height / canvasScale / 2 - typeIconSize / 4) * canvasScale,
    )
    context.scale(canvasScale, canvasScale)
    typeIcon.src = typeIcon.src.replace('currentColor', 'white')
    if (typeIcon.complete) {
      context.drawImage(typeIcon, 0, 0, typeIconSize, typeIconSize)
    }
    context.restore()

    // Fill label text
    const labelText = node.metadata.name
    context.font = `${8 * canvasScale}px Schibsted Grotesk, sans-serif`
    context.fillStyle = typeTextColor
    context.fillText(labelText, x + 5 * canvasScale, y + 13 * canvasScale)

    // Fill label icon
    const labelIconSize = 64
    context.save()
    context.translate(
      x - (width / (canvasScale * 2) - labelIconSize) * canvasScale,
      y + (height / canvasScale / 2 - labelIconSize / 2.75) * canvasScale,
    )
    context.scale(canvasScale, canvasScale)
    labelIcon.src = labelIcon.src.replace('currentColor', encodeURIComponent(labelTextColor))
    if (labelIcon.complete) {
      context.drawImage(labelIcon, 0, 0, labelIconSize, labelIconSize)
    }
    context.restore()

    // Draw node error icon
    if (nodeErrors && nodeErrors[nodeId] && nodeErrors[nodeId].length > 0 && nodeIcons.error) {
      const errorIconSize = 26
      context.save()
      context.translate(x - 42 * canvasScale, y)
      context.scale(canvasScale, canvasScale)
      nodeIcons.error.src = nodeIcons.error.src.replace(
        'currentColor',
        encodeURIComponent('#DE6237'),
      )
      context.drawImage(nodeIcons.error, 0, 0, errorIconSize, errorIconSize)
      context.restore()
    }

    // Draw input connection points
    for (const input of node.metadata.io.inputs) {
      if (input.id === 'passThrough') continue
      const inputIoPoint: IoPoint = {
        io: 'input',
        coordinates: {
          x: x + width / 2 - 2.5,
          y: y,
        },
        nodeId: nodeId,
        channelNodeIo: input,
      }

      // Check if current edge is attached to this IO point
      const isCurrentEdgeAttached = isCurrentEdgeAttachedToIoPoint(currentEdge, inputIoPoint)

      // Check if this input is connected to the selected node
      const isConnectedToSelectedNode =
        currentNode &&
        edges.size > 0 &&
        Array.from(edges.values()).some(
          (edge) =>
            edge.to.nodeId === nodeId &&
            edge.to.ioId === input.id &&
            edge.from.nodeId === currentNode.id,
        )

      // Get color based on node type
      const ioColor = (() => {
        if (
          (currentNode && nodeId === currentNode.id) ||
          isCurrentEdgeAttached ||
          isConnectedToSelectedNode
        ) {
          const currentNodeType =
            currentNode?.metadata.type || nodes.get(currentEdge!.from.nodeId)?.metadata.type

          switch (currentNodeType) {
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

      // Draw horizontal line
      const lineWidth = 20 * canvasScale
      context.beginPath()
      context.moveTo(inputIoPoint.coordinates.x - lineWidth / 2, inputIoPoint.coordinates.y)
      context.lineTo(inputIoPoint.coordinates.x + lineWidth / 2, inputIoPoint.coordinates.y)
      context.strokeStyle = ioColor
      context.lineWidth = 4 * canvasScale
      context.stroke()

      const existingIoPoint = ioPoints.get(`${node.id}_${input.id}`)
      if (
        !existingIoPoint ||
        existingIoPoint.coordinates.x !== inputIoPoint.coordinates.x ||
        existingIoPoint.coordinates.y !== inputIoPoint.coordinates.y
      ) {
        upsertIoPoint(`${node.id}_${input.id}`, inputIoPoint)
      }
    }

    const numOutputs = node.metadata.io.outputs.length
    // Draw output connection points
    for (const [index, output] of node.metadata.io.outputs.entries()) {
      // Hide HTTP response output points - these are implicitly created by the rules engine
      if (output.id === 'httpResponseOutput') continue

      // Calculate the offset based on the number of outputs and the index
      const spacing = 30 * canvasScale // Adjust this value to control the spacing between outputs
      const totalWidth = (numOutputs - 1) * spacing // Total width occupied by all outputs
      const offsetX = index * spacing - totalWidth / 2 // Center the outputs

      const outputIoPoint: IoPoint = {
        io: 'output',
        coordinates: {
          x: x + width / 2 - 2.5 + offsetX,
          y: y + height,
        },
        nodeId: nodeId,
        channelNodeIo: output,
      }

      // Check if current edge is attached to this IO point
      const isCurrentEdgeAttached = isCurrentEdgeAttachedToIoPoint(currentEdge, outputIoPoint)

      // Get color for output points (reusing same logic)
      const outputIoColor = (() => {
        if (isCurrentEdgeAttached || (currentNode && nodeId === currentNode.id)) {
          switch (node.metadata.type) {
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

      context.beginPath()
      context.ellipse(
        outputIoPoint.coordinates.x,
        outputIoPoint.coordinates.y,
        4.5 * canvasScale,
        4.5 * canvasScale,
        0,
        0,
        360,
      )
      context.fillStyle = outputIoColor
      context.fill()

      // Add labels to output IO points
      if (node.metadata.type === 'conditional') {
        context.font = `${9 * canvasScale}px Schibsted Grotesk, sans-serif`
        const outputIdTextMetrics = context.measureText(output.id)

        // Draw the text
        context.fillStyle = typeTextColor
        context.fillText(
          output.id,
          outputIoPoint.coordinates.x - outputIdTextMetrics.width / 2,
          outputIoPoint.coordinates.y - 7.5 * canvasScale,
        )
      }

      // Upsert IO point to ioPoints map
      const existingIoPoint = ioPoints.get(`${node.id}_${output.id}`)
      if (
        !existingIoPoint ||
        existingIoPoint.coordinates.x !== outputIoPoint.coordinates.x ||
        existingIoPoint.coordinates.y !== outputIoPoint.coordinates.y
      ) {
        upsertIoPoint(`${node.id}_${output.id}`, outputIoPoint)
      }
    }
  }
}

export default drawNodes
