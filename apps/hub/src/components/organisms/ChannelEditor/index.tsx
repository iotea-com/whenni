'use client'

import { FC, useEffect, useRef, useState } from 'react'
import { useDrawChannelEditor } from './useDrawChannelEditor'
import { useChannelEditorKeyboardEvents } from './events/keyboard'
import { useChannelEditorMouseEvents } from './events/mouse'
import isBrowser from '@gruent/libs/frontend/util/isBrowser'
import useChannelEditorStore from '@gruent/hub/stores/channelEditorStore'
import {
  RemixIcon,
  riCrosshair2Line,
  riErrorWarningFill,
  riZoomInLine,
  riZoomOutLine,
} from '@mwarnerdotme/react-remixicon'
import { openModal } from '@gruent/libs/frontend/hooks/useModal'
import ChannelValidationErrorModal from './ChannelValidationErrorModal'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import noop from '@gruent/libs/frontend/util/noop'

const GRID_SIZE = 24
const CANVAS_SCALE = 2

export type EditorHooks = {
  onSave?: () => void
  onExport?: () => void
  onChangeMode?: () => void
}

type Props = {
  hooks?: EditorHooks
}

const ChannelEditor: FC<Props> = ({ hooks }) => {
  // Initialize shared state
  const editorRef = useRef<HTMLCanvasElement>(null)
  const [editorSize, setEditorSize] = useState({ width: 0, height: 0 })

  // Set up store
  const validationErrors = useChannelEditorStore((state) => state.validationErrors)
  const setPan = useChannelEditorStore((state) => state.setPan)
  const zoom = useChannelEditorStore((state) => state.zoom)
  const setZoom = useChannelEditorStore((state) => state.setZoom)

  // Initialize draw canvas function
  useDrawChannelEditor(editorRef, editorSize, GRID_SIZE, CANVAS_SCALE)

  // Initialize keyboard events
  useChannelEditorKeyboardEvents(editorRef, {
    onSave: hooks?.onSave,
    onExport: hooks?.onExport,
    onChangeMode: hooks?.onChangeMode,
  })

  // Initialize mouse events
  const { handleMouseMove, handleMouseUp, handleMouseDown, handleDragOver, handleDrop } =
    useChannelEditorMouseEvents(editorRef, GRID_SIZE, CANVAS_SCALE)

  // Initialize canvas size
  useEffect(() => {
    if (!isBrowser) return
    if (!editorRef.current) return

    const containerElement = editorRef.current.parentElement

    setEditorSize({
      width: containerElement?.offsetWidth ?? 0,
      height: containerElement?.offsetHeight ?? 0,
    })
  }, [])

  // Handle page resize
  useEffect(() => {
    if (!isBrowser) return

    const handleResize = () => {
      if (!editorRef || !editorRef.current) return

      const containerElement = editorRef.current.parentElement

      setEditorSize({
        width: containerElement?.offsetWidth ?? 0,
        height: containerElement?.offsetHeight ?? 0,
      })
    }

    window.addEventListener('resize', handleResize)

    return () => window.removeEventListener('resize', handleResize)
  }, [editorRef])

  return (
    <div className="relative grow flex flex-col">
      <ChannelValidationErrorModal />
      {validationErrors &&
        validationErrors.channelErrors &&
        validationErrors.channelErrors.length > 0 && (
          <RemixIcon
            icon={riErrorWarningFill}
            size="xl"
            className="text-error absolute right-7 top-7 cursor-pointer select-none z-10"
            onClick={() => openModal('channelValidationErrors')}
          />
        )}

      <canvas
        ref={editorRef}
        tabIndex={0}
        onMouseDown={handleMouseDown}
        onMouseUp={handleMouseUp}
        onMouseOut={handleMouseUp}
        onMouseMove={handleMouseMove}
        onDragOver={handleDragOver}
        onDrop={handleDrop}
        width={editorSize.width}
        height={editorSize.height}
        className="grow focus:outline-hidden overflow-hidden relative"
      />

      <div className="absolute left-8 bottom-6 flex flex-col gap-1">
        <div className="flex gap-2">
          <Button
            variant="transparent"
            className="p-1! bg-gray-50! dark:bg-gray-900! border-gray-400! hover:border-green-600! text-gray-500! hover:text-green-600!"
            onClick={() => (zoom < 1.5 ? setZoom(zoom + 0.25) : noop)}
          >
            <RemixIcon icon={riZoomInLine} />
          </Button>
          <Button
            variant="transparent"
            className="p-1! bg-gray-50! dark:bg-gray-900! border-gray-400! hover:border-green-600! text-gray-500! hover:text-green-600!"
            onClick={() => (zoom > 0.5 ? setZoom(zoom - 0.25) : noop)}
          >
            <RemixIcon icon={riZoomOutLine} />
          </Button>
          <Button
            variant="transparent"
            className="p-1! bg-gray-50! dark:bg-gray-900! border-gray-400! hover:border-green-600! text-gray-500! hover:text-green-600!"
            onClick={() => {
              setPan(0, 0)
              setZoom(1)
            }}
          >
            <RemixIcon icon={riCrosshair2Line} />
          </Button>
        </div>
      </div>
    </div>
  )
}

export default ChannelEditor
