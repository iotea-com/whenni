import { RefObject, useCallback, useEffect, useRef, useState } from 'react'

import useChannelEditorStore from '@iotea/hub/stores/channelEditorStore'
import useSettingsStore from '@iotea/hub/stores/settingsStore'
import { drawBlankCanvas, drawGrid, drawNodes, drawEdges, drawNotes } from './draw'
import theme from '@iotea/libs/frontend/themes/tailwind'

export const useDrawChannelEditor = (
  editorRef: RefObject<HTMLCanvasElement | null>,
  editorSize: {
    width: number
    height: number
  },
  gridSize: number,
  canvasScale: number,
) => {
  // Set up store
  const nodes = useChannelEditorStore((state) => state.nodes)
  const edges = useChannelEditorStore((state) => state.edges)
  const ioPoints = useChannelEditorStore((state) => state.ioPoints)
  const currentNode = useChannelEditorStore((state) => state.currentNode)
  const currentEdge = useChannelEditorStore((state) => state.currentEdge)
  const currentNote = useChannelEditorStore((state) => state.currentNote)
  const notes = useChannelEditorStore((state) => state.notes)
  const upsertIoPoint = useChannelEditorStore((state) => state.upsertIoPoint)
  const pan = useChannelEditorStore((state) => state.pan)
  const userTheme = useSettingsStore((state) => state.theme)
  const showChannelEditorGrid = useSettingsStore((state) => state.showChannelEditorGrid)
  const zoom = useChannelEditorStore((state) => state.zoom)
  const edgePreview = useChannelEditorStore((state) => state.edgePreview)
  const validationErrors = useChannelEditorStore((state) => state.validationErrors)

  const [dashOffset, _setDashOffset] = useState(0) // causes performance issues

  const drawCanvas = useCallback((): RefObject<HTMLCanvasElement | null> | undefined => {
    const editor = editorRef.current
    if (!editorRef || !editor) return editorRef
    const context = editor.getContext('2d')
    if (!context) return

    // Return if SSR
    if (editorSize.width <= 0 || editorSize.height <= 0) return

    // setDashOffset((prev) => (prev + 0.5)) // causes performance issues

    // Adjust canvas resolution
    editor.style.width = editorSize.width + 'px'
    editor.style.height = editorSize.height + 'px'
    editor.width = editorSize.width * canvasScale
    editor.height = editorSize.height * canvasScale
    const { width, height } = editor

    // Reset the canvas
    drawBlankCanvas({ context, width, height, zoom, editor, userTheme })

    // Draw grid
    if (showChannelEditorGrid)
      drawGrid({ context, width, height, gridSize, canvasScale, pan, userTheme })

    // Draw edges
    drawEdges({
      context,
      canvasScale,
      nodes,
      edges,
      ioPoints,
      currentNode,
      currentEdge,
      dashOffset,
      width,
      height,
      zoom,
      editor,
      userTheme,
    })

    // Draw nodes
    drawNodes({
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
      width,
      height,
      zoom,
      editor,
      nodeErrors: validationErrors ? validationErrors.nodeErrors : null,
    })

    // Draw notes
    drawNotes({
      context,
      notes,
      canvasScale,
      pan,
      userTheme,
      currentNote,
      width,
      height,
      zoom,
      editor,
    })

    // Draw edge preview
    if (edgePreview) {
      context.beginPath()
      context.moveTo(edgePreview.from.x, edgePreview.from.y)
      context.lineTo(edgePreview.to.x * canvasScale, edgePreview.to.y * canvasScale)
      context.strokeStyle = theme.colors.primary
      context.lineWidth = 1
      context.stroke()
    }

    // Restore context transformation
    context.restore()

    return editorRef
  }, [
    nodes,
    edges,
    notes,
    editorRef,
    upsertIoPoint,
    ioPoints,
    currentNode,
    currentEdge,
    currentNote,
    showChannelEditorGrid,
    gridSize,
    editorSize,
    canvasScale,
    dashOffset,
    pan,
    userTheme,
    zoom,
    edgePreview,
    validationErrors,
  ])

  // Re-render the canvas when the drawCanvas function is updated - uses requestAnimationFrame to
  // limit the rate at which the canvas is drawn
  const requestRef = useRef<number>(0)

  useEffect(() => {
    requestRef.current = requestAnimationFrame(drawCanvas)
    return () => {
      if (requestRef.current) cancelAnimationFrame(requestRef.current)
    }
  }, [drawCanvas])
}
