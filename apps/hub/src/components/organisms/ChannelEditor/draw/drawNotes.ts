import { ChannelNote } from '@iotea/libs/engine/channels/channels'
import { calculateNoteDisplay } from '../util/calculateNoteDisplay'
import { Theme } from '@iotea/hub/stores/settingsStore'
import { wrapText } from '../util/wrapText'

type DrawNotesParams = {
  context: CanvasRenderingContext2D
  width: number
  height: number
  zoom: number
  editor: HTMLCanvasElement
  notes: Map<string, ChannelNote>
  currentNote: ChannelNote | null
  canvasScale: number
  pan: { x: number; y: number }
  userTheme: Theme
}

const drawNotes = ({
  context,
  notes,
  canvasScale,
  pan,
  userTheme,
  currentNote,
}: DrawNotesParams) => {
  for (const noteEntry of notes) {
    const noteId = noteEntry[0]
    const note = noteEntry[1]

    const {
      x,
      y,
      width,
      height,
      backgroundColor,
      borderColor,
      borderWidth,
      labelTextColor,
      typeTextColor,
    } = calculateNoteDisplay(note, canvasScale, pan, userTheme)

    // fill rectangle background
    context.beginPath()
    context.rect(x, y, width, height)
    context.fillStyle = backgroundColor
    context.fill()

    // fill type text
    context.font = `${9 * canvasScale}px Schibsted Grotesk, sans-serif`
    context.fillStyle = typeTextColor
    context.fillText('NOTE', x + 10 * canvasScale, y + 15 * canvasScale)

    // wrap label text to prevent overflow
    context.font = `${14 * canvasScale}px Schibsted Grotesk, sans-serif`
    context.fillStyle = labelTextColor

    const lineHeight = 20 * canvasScale
    const maxTextWidth = width - 20 * canvasScale
    const lines = wrapText(context, note.text, maxTextWidth)

    let textY = y + 32.5 * canvasScale
    for (const line of lines) {
      // Prevent vertical overflow
      if (textY + lineHeight > y + height) break

      context.fillText(line, x + 10 * canvasScale, textY)
      textY += lineHeight
    }

    // apply stroke if this is the current note
    if (currentNote && noteId === currentNote.id) {
      context.strokeStyle = borderColor
      context.lineWidth = borderWidth * 2 * canvasScale
      context.stroke()
    } else {
      context.strokeStyle = borderColor
      context.lineWidth = borderWidth * canvasScale
      context.stroke()
    }
  }
}

export default drawNotes
