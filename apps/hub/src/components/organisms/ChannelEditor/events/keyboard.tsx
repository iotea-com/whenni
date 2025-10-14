'use client'

import { useKeyDown } from '@react-hooks-library/core'
import { RefObject } from 'react'
import useChannelEditorStore from '@iotea/hub/stores/channelEditorStore'
import { EditorHooks } from '..'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'

export const useChannelEditorKeyboardEvents = (
  editorRef: RefObject<HTMLCanvasElement | null>,
  hooks?: EditorHooks,
) => {
  // set up store
  const currentNode = useChannelEditorStore((state) => state.currentNode)
  const removeNode = useChannelEditorStore((state) => state.removeNode)
  const setCurrentNode = useChannelEditorStore((state) => state.setCurrentNode)
  const edges = useChannelEditorStore((state) => state.edges)
  const removeEdge = useChannelEditorStore((state) => state.removeEdge)
  const ioPoints = useChannelEditorStore((state) => state.ioPoints)
  const removeIoPoint = useChannelEditorStore((state) => state.removeIoPoint)
  const currentEdge = useChannelEditorStore((state) => state.currentEdge)
  const setCurrentEdge = useChannelEditorStore((state) => state.setCurrentEdge)
  const currentNote = useChannelEditorStore((state) => state.currentNote)
  const setCurrentNote = useChannelEditorStore((state) => state.setCurrentNote)
  const removeNote = useChannelEditorStore((state) => state.removeNote)

  const handleDeleteKeypress = () => {
    if (!editorRef || !editorRef.current) return
    if (document.activeElement !== editorRef.current) return

    if (currentNode) {
      // update nodes map
      const currentNodeEdges = Array.from(edges.values()).map((edge) => {
        if (edge.from.nodeId === currentNode.id) return [edge.id, true]
        if (edge.to.nodeId === currentNode.id) return [edge.id, true]
        return [edge.id, false]
      })

      for (const currentNodeEdge of currentNodeEdges) {
        const [edgeId, isAttachedToEdge] = currentNodeEdge
        if (isAttachedToEdge) removeEdge(edgeId as string)
      }

      // update IO points map
      for (const [id, _] of ioPoints) {
        if (id.includes(currentNode.id)) removeIoPoint(id)
      }

      // update nodes map
      removeNode(currentNode.id)

      // reset current node
      setCurrentNode(null)
    }

    if (currentNote) {
      // update notes map
      removeNote(currentNote.id)

      // reset current edge
      setCurrentNote(null)
    }

    if (currentEdge) {
      // update edges map
      removeEdge(currentEdge.id)

      // reset current edge
      setCurrentEdge(null)
    }
  }

  const handleSaveKeypress = (e: KeyboardEvent) => {
    if ((e.ctrlKey || e.metaKey) && e.key === 's') {
      e.preventDefault()
      hooks?.onSave?.()
    }
  }

  const handleExportKeypress = (e: KeyboardEvent) => {
    if ((e.ctrlKey || e.metaKey) && e.key === 'e') {
      e.preventDefault()
      hooks?.onExport?.()
    }
  }

  const handleModeKeypress = (e: KeyboardEvent) => {
    if ((e.ctrlKey || e.metaKey) && e.key === 'i') {
      e.preventDefault()
      hooks?.onChangeMode?.()
    }
  }

  const handleUndoKeypress = (e: KeyboardEvent) => {
    if (document.activeElement !== editorRef.current) return

    // Handle redo (shift)
    if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.key === 'z') {
      e.preventDefault()
      addToast({
        title: 'Coming soon',
        body: 'Redo is coming soon!',
        level: 'info',
      })
      return
    }

    // Handle undo (no shift)
    if ((e.ctrlKey || e.metaKey) && e.key === 'z') {
      e.preventDefault()
      addToast({
        title: 'Coming soon',
        body: 'Undo is coming soon!',
        level: 'info',
      })
    }
  }

  const handleCopyKeypress = (e: KeyboardEvent) => {
    if (document.activeElement !== editorRef.current) return

    // Handle copy
    if ((e.ctrlKey || e.metaKey) && e.key === 'c') {
      e.preventDefault()
      addToast({
        title: 'Coming soon',
        body: 'Copying nodes is coming soon!',
        level: 'info',
      })
    }
  }

  const handlePasteKeypress = (e: KeyboardEvent) => {
    if (document.activeElement !== editorRef.current) return

    // Handle paste
    if ((e.ctrlKey || e.metaKey) && e.key === 'v') {
      e.preventDefault()
      addToast({
        title: 'Coming soon',
        body: 'Pasting nodes is coming soon!',
        level: 'info',
      })
    }
  }

  const handleDuplicateKeypress = (e: KeyboardEvent) => {
    // Handle paste
    if ((e.ctrlKey || e.metaKey) && e.key === 'd') {
      e.preventDefault()
      addToast({
        title: 'Coming soon',
        body: 'Duplicating nodes is coming soon!',
        level: 'info',
      })
    }
  }

  useKeyDown(['Delete', 'Backspace'], handleDeleteKeypress)
  useKeyDown(['s'], handleSaveKeypress)
  useKeyDown(['e'], handleExportKeypress)
  useKeyDown(['i'], handleModeKeypress)
  useKeyDown(['z'], handleUndoKeypress)
  useKeyDown(['c'], handleCopyKeypress)
  useKeyDown(['v'], handlePasteKeypress)
  useKeyDown(['d'], handleDuplicateKeypress)

  return { handleDeleteKeypress }
}
