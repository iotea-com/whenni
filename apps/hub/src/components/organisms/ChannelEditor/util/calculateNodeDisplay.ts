import { ChannelNode } from '@iotea/libs/engine/nodes/v1'

import theme from '@iotea/libs/frontend/themes/tailwind'

import nodeIcons, { cloneImage } from './nodeIcons'

type NodeDisplay = {
  x: number
  y: number
  width: number
  height: number
  backgroundColor: string
  borderColor: string
  borderWidth: 1 | 2 | 3 | 4
  selectedBorderColor: string
  selectedBorderWidth: 1 | 2 | 3 | 4
  accentColor: string
  typeTextColor: string
  labelTextColor: string
  typeIcon: HTMLImageElement
  labelIcon: HTMLImageElement
}

export const calculateNodeDisplay = (
  node: ChannelNode,
  canvasScale: number,
  pan: { x: number; y: number } = { x: 0, y: 0 },
  userTheme: 'light' | 'dark' = 'light',
) => {
  const nodeDisplay: NodeDisplay = (() => {
    const base: NodeDisplay = {
      x: node.coordinates.x * canvasScale - pan.x,
      y: node.coordinates.y * canvasScale - pan.y,
      width: 90 * canvasScale,
      height: 90 * canvasScale,
      backgroundColor: userTheme === 'light' ? theme.colors.gray[50] : theme.colors.gray[900],
      borderColor: userTheme === 'light' ? theme.colors.gray[300] : theme.colors.gray[400],
      borderWidth: 1,
      selectedBorderColor: userTheme === 'light' ? theme.colors.gray[600] : theme.colors.gray[300],
      selectedBorderWidth: 1,
      accentColor: userTheme === 'light' ? theme.colors.gray[100] : theme.colors.gray[800],
      typeTextColor: userTheme === 'light' ? theme.colors.gray[500] : theme.colors.gray[400],
      labelTextColor: userTheme === 'light' ? theme.colors.gray[800] : theme.colors.gray[200],
      typeIcon: nodeIcons.question ?? new Image(),
      labelIcon: nodeIcons.question ?? new Image(),
    }

    const labelIcon = calculateNodeLabelIcon(node.metadata.label)

    switch (node.metadata.type) {
      case 'source':
        return {
          ...base,
          selectedBorderColor: theme.colors.green[500],
          borderColor: theme.colors.green[200],
          accentColor: theme.colors.green[400],
          typeTextColor: userTheme === 'light' ? theme.colors.green[700] : theme.colors.gray[50],
          labelTextColor: userTheme === 'light' ? theme.colors.green[800] : theme.colors.white,
          typeIcon: cloneImage(nodeIcons.source ?? new Image()),
          labelIcon,
        }
      case 'processing':
        return {
          ...base,
          selectedBorderColor: theme.colors.blue[500],
          borderColor: theme.colors.blue[200],
          accentColor: theme.colors.blue[400],
          typeTextColor: userTheme === 'light' ? theme.colors.blue[700] : theme.colors.gray[50],
          typeIcon: cloneImage(nodeIcons.processing ?? new Image()),
          labelTextColor: userTheme === 'light' ? theme.colors.blue[900] : theme.colors.white,
          labelIcon,
        }
      case 'action':
        return {
          ...base,
          selectedBorderColor: theme.colors.red[500],
          borderColor: theme.colors.red[200],
          accentColor: theme.colors.red[400],
          typeTextColor: userTheme === 'light' ? theme.colors.red[700] : theme.colors.gray[50],
          typeIcon: cloneImage(nodeIcons.action ?? new Image()),
          labelTextColor: userTheme === 'light' ? theme.colors.red[800] : theme.colors.white,
          labelIcon,
        }
      case 'conditional':
        return {
          ...base,
          selectedBorderColor: theme.colors.orange[500],
          borderColor: theme.colors.orange[200],
          accentColor: theme.colors.orange[400],
          typeTextColor: userTheme === 'light' ? theme.colors.orange[700] : theme.colors.gray[50],
          typeIcon: cloneImage(nodeIcons.conditional ?? new Image()),
          labelTextColor: userTheme === 'light' ? theme.colors.orange[900] : theme.colors.white,
          labelIcon,
        }
      case 'custom':
        return {
          ...base,
          selectedBorderColor: theme.colors.violet[500],
          borderColor: theme.colors.violet[200],
          accentColor: theme.colors.violet[400],
          typeTextColor: userTheme === 'light' ? theme.colors.violet[700] : theme.colors.gray[50],
          labelTextColor: userTheme === 'light' ? theme.colors.violet[900] : theme.colors.white,
          typeIcon: cloneImage(nodeIcons.custom ?? new Image()),
          labelIcon,
        }
      default:
        return base
    }
  })()

  return nodeDisplay
}

const calculateNodeLabelIcon = (label: string) => {
  switch (label) {
    case 'mqtt':
    case 'messageQueue':
      return cloneImage(nodeIcons.messageQueue ?? new Image())
    case 'http':
    case 'httpResponse':
      return cloneImage(nodeIcons.http ?? new Image())
    case 'timer':
      return cloneImage(nodeIcons.timer ?? new Image())
    case 'fileStorage':
      return cloneImage(nodeIcons.fileStorage ?? new Image())
    case 'documentDb':
    case 'timeSeriesDb':
      return cloneImage(nodeIcons.database ?? new Image())
    case 'log':
      return cloneImage(nodeIcons.log ?? new Image())
    case 'metric':
      return cloneImage(nodeIcons.metric ?? new Image())
    case 'notification':
      return cloneImage(nodeIcons.notification ?? new Image())
    case 'boolean':
    case 'stringCompare':
    case 'threshold':
    case 'existence':
      return cloneImage(nodeIcons.flowChart ?? new Image())
    case 'transform':
      return cloneImage(nodeIcons.shapes ?? new Image())
    default:
      return cloneImage(nodeIcons.question ?? new Image())
  }
}
