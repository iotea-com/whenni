import { ChannelNote } from '@gruent/libs/engine/channels/index'

import theme from '@gruent/libs/frontend/themes/tailwind'

type NoteDisplay = {
  x: number
  y: number
  width: number
  height: number
  backgroundColor: string
  borderColor: string
  borderWidth: 1 | 2 | 3 | 4
  typeTextColor: string
  labelTextColor: string
}

export const calculateNoteDisplay = (
  note: ChannelNote,
  canvasScale: number,
  pan: { x: number; y: number } = { x: 0, y: 0 },
  userTheme: 'light' | 'dark' = 'light',
) => {
  const noteDisplay: NoteDisplay = (() => {
    return {
      x: note.coordinates.x * canvasScale - pan.x,
      y: note.coordinates.y * canvasScale - pan.y,
      width: 192 * canvasScale,
      height: 95 * canvasScale,
      backgroundColor: userTheme === 'light' ? theme.colors.gray[50] : theme.colors.gray[900],
      borderColor: userTheme === 'light' ? theme.colors.gray[700] : theme.colors.gray[100],
      borderWidth: 1,
      typeTextColor: userTheme === 'light' ? theme.colors.gray[500] : theme.colors.gray[400],
      labelTextColor: userTheme === 'light' ? theme.colors.gray[800] : theme.colors.gray[200],
    }
  })()

  return noteDisplay
}
