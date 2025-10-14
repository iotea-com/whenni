import { ChannelNode } from '@iotea/libs/engine/nodes/v1'
import { ChannelEdge, ChannelEdgeConnection, ChannelNote } from '@iotea/libs/engine/channels/index'

import { DragEventHandler, MouseEventHandler, RefObject, useCallback, useState } from 'react'
import { isMouseInIoPoint } from '../util/isMouseInIoPoint'
import { isMouseInNode } from '../util/isMouseInNode'
import { IoPoint } from '../types/IoPoint'

import { isMouseOnEdge } from '../util/isMouseOnEdge'
import useChannelEditorStore from '@iotea/hub/stores/channelEditorStore'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import { isMouseInNote } from '../util/isMouseInNote'
import { generateEdgeId, generateNodeId, generateNoteId } from '@iotea/hub/util/idGenerator'

export const useChannelEditorMouseEvents = (
  editorRef: RefObject<HTMLCanvasElement | null>,
  gridSize: number,
  canvasScale: number,
) => {
  const [isCreatingEdge, setIsCreatingEdge] = useState<IoPoint | null>(null)
  const [dragging, setDragging] = useState<boolean>(false)
  const [diffMouseX, setDiffMouseX] = useState<number>(0)
  const [diffMouseY, setDiffMouseY] = useState<number>(0)

  // set up store
  const nodes = useChannelEditorStore((state) => state.nodes)
  const upsertNode = useChannelEditorStore((state) => state.upsertNode)
  const edges = useChannelEditorStore((state) => state.edges)
  const upsertEdge = useChannelEditorStore((state) => state.upsertEdge)
  const notes = useChannelEditorStore((state) => state.notes)
  const upsertNote = useChannelEditorStore((state) => state.upsertNote)
  const ioPoints = useChannelEditorStore((state) => state.ioPoints)
  const currentNode = useChannelEditorStore((state) => state.currentNode)
  const currentNote = useChannelEditorStore((state) => state.currentNote)
  const setCurrentNode = useChannelEditorStore((state) => state.setCurrentNode)
  const setCurrentEdge = useChannelEditorStore((state) => state.setCurrentEdge)
  const setCurrentNote = useChannelEditorStore((state) => state.setCurrentNote)
  const pan = useChannelEditorStore((state) => state.pan)
  const setPan = useChannelEditorStore((state) => state.setPan)
  const setEdgePreview = useChannelEditorStore((state) => state.setEdgePreview)

  const handleMouseDown: MouseEventHandler<HTMLCanvasElement> = useCallback(
    (e) => {
      e.preventDefault()
      editorRef.current?.focus() // add focus so that keyboard events only apply when the canvas is focused

      const editor = editorRef.current
      if (!editor) return

      const canvasRect = editor.getBoundingClientRect()
      const mouseX = e.clientX - canvasRect.left
      const mouseY = e.clientY - canvasRect.top

      setIsCreatingEdge(null)

      // if attempting to drag an ioPoint, prepare to update edges
      for (const [_, ioPoint] of ioPoints) {
        const inIoPoint = isMouseInIoPoint(mouseX, mouseY, ioPoint, canvasScale)
        if (inIoPoint) {
          setIsCreatingEdge(ioPoint)
          return
        }
      }

      // if attempting to select an edge
      const newCurrentEdge = (() => {
        for (const [_, edge] of edges) {
          const hoveredEdge = isMouseOnEdge(
            mouseX,
            mouseY,
            edge,
            Array.from(ioPoints.values()),
            canvasScale,
          )

          if (hoveredEdge) return hoveredEdge
        }
      })()

      setCurrentEdge(newCurrentEdge ?? null)
      if (newCurrentEdge) return

      // if attempting to select a note, update currentNote
      const newCurrentNote = (() => {
        for (const note of notes.entries()) {
          const noteId = note[0]
          const noteValue = note[1]
          if (isMouseInNote(mouseX, mouseY, noteValue, canvasScale, pan))
            return { noteId, noteValue }
        }
      })()

      setCurrentNote(newCurrentNote?.noteValue ?? null)
      if (newCurrentNote) {
        // set current note
        setDragging(true)

        setDiffMouseX(mouseX - newCurrentNote.noteValue.coordinates.x)
        setDiffMouseY(mouseY - newCurrentNote.noteValue.coordinates.y)

        // bring it to the front of the z dimension
        upsertNote(newCurrentNote.noteId, newCurrentNote.noteValue)
        return
      }

      // if attempting to select a node, update currentNode
      const newCurrentNode = (() => {
        for (const node of nodes.entries()) {
          const nodeId = node[0]
          const nodeValue = node[1]
          if (isMouseInNode(mouseX, mouseY, nodeValue, canvasScale, pan))
            return { nodeId, nodeValue }
        }
      })()

      setCurrentNode(newCurrentNode?.nodeValue ?? null)
      if (newCurrentNode) {
        // set current node
        setDragging(true)

        setDiffMouseX(mouseX - newCurrentNode.nodeValue.coordinates.x)
        setDiffMouseY(mouseY - newCurrentNode.nodeValue.coordinates.y)

        // bring it to the front of the z dimension
        upsertNode(newCurrentNode.nodeId, newCurrentNode.nodeValue)
        return
      }

      // Reset if no nodes, edges, or notes are clicked or dragged
      setDragging(true)
      // Save the initial mouse position
      setDiffMouseX(e.clientX)
      setDiffMouseY(e.clientY)
    },
    [
      nodes,
      edges,
      notes,
      ioPoints,
      editorRef,
      setCurrentNode,
      setCurrentEdge,
      setCurrentNote,
      canvasScale,
      upsertNode,
      upsertNote,
      pan,
    ],
  )

  const handleMouseUp: MouseEventHandler<HTMLCanvasElement> = useCallback(
    (e) => {
      e.preventDefault()

      const editor = editorRef.current
      if (!editor) return

      const canvasRect = editor.getBoundingClientRect()
      const mouseX = e.clientX - canvasRect.left
      const mouseY = e.clientY - canvasRect.top

      setDragging(false)
      const mouseCursor = calculateMouseCursor({
        currentNode,
        currentNote,
        nodes,
        notes,
        edges,
        ioPoints,
        mouseX,
        mouseY,
        dragging: false,
        isCreatingEdge,
        canvasScale,
        pan,
      })
      editor.style.cursor = mouseCursor

      if (isCreatingEdge) {
        for (const [_, ioPoint] of ioPoints) {
          if (isMouseInIoPoint(mouseX, mouseY, ioPoint, canvasScale)) {
            // ignore the connection if it's to the same node
            if (isCreatingEdge.nodeId === ioPoint.nodeId) {
              addToast({
                title: 'Could not create edge',
                body: 'Cannot connect a node to itself.',
                level: 'warning',
              })
              break
            }

            // ignore the connection if start io and end io are the same
            const startIo = isCreatingEdge.io
            const endIo = ioPoint.io
            if (startIo === endIo) {
              addToast({
                title: 'Could not create edge',
                body: `Cannot connect an ${startIo} to another ${endIo}.`,
                level: 'warning',
              })
              break
            }

            // a
            const outputNodeIoPoint = (() => {
              if (isCreatingEdge.io === 'output') return isCreatingEdge
              return ioPoint
            })()

            const inputNodeIoPoint = (() => {
              if (isCreatingEdge.io === 'input') return isCreatingEdge
              return ioPoint
            })()

            const inputNode = nodes.get(inputNodeIoPoint.nodeId)
            const outputNode = nodes.get(outputNodeIoPoint.nodeId)
            if (!inputNode || !outputNode) {
              // TODO: toast notification
              break
            }

            // TODO: ignore the connection if there is already an edge for these two nodes

            // TODO: ignore the connection if it's not in the available connections array
            const isAllowedInput: boolean = (() => {
              const outputNodeType = outputNode.metadata.type
              for (const io of inputNode.metadata.io.inputs) {
                if (io.allowedEdge.includes(outputNodeType) || io.allowedEdge.includes('*'))
                  return true
              }

              return false
            })()

            const isAllowedOutput: boolean = (() => {
              const inputNodeType = inputNode.metadata.type
              for (const io of outputNode.metadata.io.outputs) {
                if (io.allowedEdge.includes(inputNodeType) || io.allowedEdge.includes('*'))
                  return true
              }

              return false
            })()

            if (!isAllowedInput || !isAllowedOutput) {
              addToast({
                title: 'Could not create edge',
                body: 'These two nodes cannot be connected.',
                level: 'warning',
              })
              break
            }

            // create a new edge in the channel config
            const from: ChannelEdgeConnection = {
              nodeId: outputNodeIoPoint.nodeId,
              ioId: outputNodeIoPoint.channelNodeIo.id,
            }

            const to: ChannelEdgeConnection = {
              nodeId: inputNodeIoPoint.nodeId,
              ioId: inputNodeIoPoint.channelNodeIo.id,
            }

            const edgeId = generateEdgeId()
            const edge: ChannelEdge = {
              id: edgeId,
              from,
              to,
            }

            upsertEdge(edgeId, edge)

            break
          }
        }

        setEdgePreview(null)
        setIsCreatingEdge(null)
      }
    },
    [
      nodes,
      isCreatingEdge,
      ioPoints,
      upsertEdge,
      editorRef,
      canvasScale,
      currentNode,
      currentNote,
      edges,
      notes,
      pan,
      setEdgePreview,
    ],
  )

  const handleMouseMove: MouseEventHandler<HTMLCanvasElement> = useCallback(
    (e) => {
      const editor = editorRef.current
      if (!editor) return

      const canvasRect = editor.getBoundingClientRect()
      const mouseX = e.clientX - canvasRect.left
      const mouseY = e.clientY - canvasRect.top

      if (currentNode && dragging) {
        // snap to the nearest grid position
        const mouseXWithDiff = mouseX - diffMouseX
        const mouseYWithDiff = mouseY - diffMouseY
        const snapX = Math.round(mouseXWithDiff / gridSize) * gridSize
        const snapY = Math.round(mouseYWithDiff / gridSize) * gridSize

        const updatedNode = {
          ...currentNode,
          coordinates: {
            x: snapX,
            y: snapY,
          },
        }

        upsertNode(currentNode.id, updatedNode)
      } else if (currentNote && dragging) {
        // snap to the nearest grid position
        const mouseXWithDiff = mouseX - diffMouseX
        const mouseYWithDiff = mouseY - diffMouseY
        const snapX = Math.round(mouseXWithDiff / gridSize) * gridSize
        const snapY = Math.round(mouseYWithDiff / gridSize) * gridSize

        const updatedNote = {
          ...currentNote,
          coordinates: {
            x: snapX,
            y: snapY,
          },
        }

        upsertNote(currentNote.id, updatedNote)
      } else if (isCreatingEdge) {
        const context = editor.getContext('2d')
        if (!context) return

        setEdgePreview({
          from: {
            x: isCreatingEdge.coordinates.x,
            y: isCreatingEdge.coordinates.y,
          },
          to: {
            x: mouseX,
            y: mouseY,
          },
        })
      } else if (dragging) {
        // Calculate the difference in mouse position
        const diffX = e.clientX - diffMouseX
        const diffY = e.clientY - diffMouseY

        setPan(pan.x - diffX, pan.y - diffY)

        // Update the saved mouse position
        setDiffMouseX(e.clientX)
        setDiffMouseY(e.clientY)
      }

      const mouseCursor = calculateMouseCursor({
        currentNode,
        currentNote,
        nodes,
        notes,
        edges,
        ioPoints,
        mouseX,
        mouseY,
        dragging,
        isCreatingEdge,
        canvasScale,
        pan,
      })
      editor.style.cursor = mouseCursor
    },
    [
      nodes,
      currentNode,
      upsertNode,
      edges,
      notes,
      currentNote,
      upsertNote,
      dragging,
      isCreatingEdge,
      ioPoints,
      gridSize,
      editorRef,
      canvasScale,
      diffMouseX,
      diffMouseY,
      pan,
      setPan,
    ],
  )

  const handleDragOver: DragEventHandler<HTMLCanvasElement> = (e) => {
    e.preventDefault()
    e.dataTransfer.dropEffect = 'copy'
    setCurrentNode(null) // reset the currentNode to ensure the OptionsPane state resets
    setCurrentEdge(null) // reset the currentEdge
  }

  const handleDrop: DragEventHandler<HTMLCanvasElement> = useCallback(
    (e) => {
      e.preventDefault()

      const editor = editorRef.current
      if (!editor) return

      const dropData = e.dataTransfer.getData('application/json')
      const parsedDropData = JSON.parse(dropData)

      if (typeof parsedDropData === 'object' && Object.keys(parsedDropData).includes('metadata')) {
        const newNode = JSON.parse(dropData) as ChannelNode
        newNode.id = generateNodeId()

        const canvasRect = editor.getBoundingClientRect()
        const mouseX = e.clientX - canvasRect.left
        const mouseY = e.clientY - canvasRect.top

        // snap to the nearest grid position
        const snapX = Math.round(mouseX / gridSize) * gridSize
        const snapY = Math.round(mouseY / gridSize) * gridSize

        newNode.coordinates.x = snapX
        newNode.coordinates.y = snapY

        upsertNode(newNode.id, newNode)
        setCurrentNode(newNode)
      }

      if (
        typeof parsedDropData === 'object' &&
        Object.keys(parsedDropData).includes('id') &&
        parsedDropData.id === '__IOTEA_NOTE__'
      ) {
        const newNote = JSON.parse(dropData) as ChannelNote
        newNote.id = generateNoteId()

        const canvasRect = editor.getBoundingClientRect()
        const mouseX = e.clientX - canvasRect.left
        const mouseY = e.clientY - canvasRect.top

        // snap to the nearest grid position
        const snapX = Math.round(mouseX / gridSize) * gridSize
        const snapY = Math.round(mouseY / gridSize) * gridSize

        newNote.coordinates.x = snapX
        newNote.coordinates.y = snapY

        upsertNote(newNote.id, newNote)
        setCurrentNote(newNote)
      }
    },
    [editorRef, upsertNode, gridSize, setCurrentNode, upsertNote, setCurrentNote],
  )

  return {
    handleMouseMove,
    handleMouseUp,
    handleMouseDown,
    handleDrop,
    handleDragOver,
  }
}

const calculateMouseCursor = ({
  currentNode,
  currentNote,
  nodes,
  notes,
  edges,
  ioPoints,
  mouseX,
  mouseY,
  dragging,
  isCreatingEdge,
  canvasScale,
  pan,
}) => {
  if ((currentNode || currentNote) && dragging) return 'move'
  if (!currentNode && !currentNote && dragging) return 'grabbing'

  for (const [_, ioPoint] of ioPoints) {
    const mouseIsInIoPoint = isMouseInIoPoint(mouseX, mouseY, ioPoint, canvasScale)
    if (mouseIsInIoPoint) return 'crosshair'
  }

  if (isCreatingEdge) return 'grabbing'

  for (const [_, node] of nodes) {
    const mouseIsInNode = isMouseInNode(mouseX, mouseY, node, canvasScale, pan)
    if (mouseIsInNode) return 'pointer'
  }

  for (const [_, note] of notes) {
    const mouseIsInNote = isMouseInNote(mouseX, mouseY, note, canvasScale, pan)
    if (mouseIsInNote) return 'pointer'
  }

  for (const [_, edge] of edges) {
    const hoveredEdge = isMouseOnEdge(
      mouseX,
      mouseY,
      edge,
      Array.from(ioPoints.values()),
      canvasScale,
    )

    if (hoveredEdge) return 'pointer'
  }

  return 'default'
}
