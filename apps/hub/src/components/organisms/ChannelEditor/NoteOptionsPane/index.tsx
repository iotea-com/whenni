import useChannelEditorStore from '@iotea/hub/stores/channelEditorStore'
import { FC, useCallback } from 'react'

type Props = {}

const NoteOptionsPane: FC<Props> = () => {
  // set up store
  const currentNote = useChannelEditorStore((state) => state.currentNote)
  const upsertNote = useChannelEditorStore((state) => state.upsertNote)
  const setCurrentNote = useChannelEditorStore((state) => state.setCurrentNote)

  const handleUpdateNote = useCallback(
    (text: string) => {
      if (!currentNote) return

      const updatedNote = {
        ...currentNote,
        text,
      }

      // Update note
      upsertNote(currentNote.id, updatedNote)

      // Update current note
      setCurrentNote(updatedNote)
    },
    [currentNote, upsertNote, setCurrentNote],
  )

  return (
    <div className="p-4 text-sm border-b">
      <h2 className="uppercase text-gray-500 font-bold text-xs mb-2">Note</h2>
      <textarea
        className="p-1 w-full min-h-24 border border-gray-200 rounded-sm"
        value={currentNote?.text}
        onChange={(e) => handleUpdateNote(e.target.value)}
      />
    </div>
  )
}

export default NoteOptionsPane
