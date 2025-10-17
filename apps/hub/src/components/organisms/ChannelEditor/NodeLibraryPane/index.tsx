import { DragEvent, FC, useMemo, useState } from 'react'
import NodeLibrarySection from './NodeLibrarySection'
import { ChannelNodeType, defaultNodes } from '@iotea/libs/engine/nodes/v1'
import {
  RemixIcon,
  riCloseFill,
  riFilePaperLine,
  riSearch2Line,
} from '@mwarnerdotme/react-remixicon'

const NodeLibraryPane: FC = () => {
  const [searchQuery, setSearchQuery] = useState<string>()
  const [openSections, setOpenSections] = useState<ChannelNodeType[]>([])

  const handleTitleClick = (nodeType: ChannelNodeType) => {
    setOpenSections((current) => {
      if (current.includes(nodeType)) {
        return current.filter((nt) => {
          if (nt === nodeType) return false
          return true
        })
      } else {
        return [...current, nodeType]
      }
    })
  }

  const handleNoteDragStart = (e: DragEvent<HTMLDivElement>) => {
    const defaultNoteValue = {
      id: '__IOTEA_NOTE__',
      coordinates: {
        x: 0,
        y: 0,
      },
      text: '',
    }
    e.dataTransfer.setData('application/json', JSON.stringify(defaultNoteValue))
  }

  const availableNodes = useMemo(() => {
    if (!searchQuery) return Array.from(defaultNodes.entries())

    return Array.from(defaultNodes.entries()).filter(([nodeName]) => {
      if (nodeName.toLowerCase().includes(searchQuery.toLowerCase())) return true
      return false
    })
  }, [searchQuery])

  return (
    <div className="flex flex-col w-full min-w-56 overflow-y-scroll overflow-x-hidden">
      <h2 className="uppercase text-gray-500 dark:text-gray-100 font-bold text-xs mb-2 pt-px">
        Nodes
      </h2>
      <div className="flex relative w-full items-bottom mb-2">
        <RemixIcon icon={riSearch2Line} className="mr-1 text-gray-700" size="sm" />
        <input
          name="nodeSearch"
          type="text"
          autoComplete="off"
          autoCorrect="off"
          spellCheck="false"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="outline-hidden border-0 border-b bg-transparent transition border-gray-400 hover:border-gray-800 focus:border-gray-800 text-xs mt-0 grow"
        />
        {searchQuery && searchQuery.length > 0 && (
          <RemixIcon
            onClick={() => setSearchQuery('')}
            className="absolute right-0 t-1/2 cursor-pointer"
            icon={riCloseFill}
          />
        )}
      </div>
      <NodeLibrarySection
        availableNodes={availableNodes}
        nodeType="source"
        open={openSections.includes('source') || availableNodes.length <= 5}
        onTitleClick={() => handleTitleClick('source')}
      />
      <NodeLibrarySection
        availableNodes={availableNodes}
        nodeType="processing"
        open={openSections.includes('processing') || availableNodes.length <= 5}
        onTitleClick={() => handleTitleClick('processing')}
      />
      <NodeLibrarySection
        availableNodes={availableNodes}
        nodeType="conditional"
        open={openSections.includes('conditional') || availableNodes.length <= 5}
        onTitleClick={() => handleTitleClick('conditional')}
      />
      <NodeLibrarySection
        availableNodes={availableNodes}
        nodeType="action"
        open={openSections.includes('action') || availableNodes.length <= 5}
        onTitleClick={() => handleTitleClick('action')}
      />
      <div className="grow" />
      <hr className="my-2" />
      <section>
        <div
          draggable
          className="my-1 px-1 select-none dark:text-gray-200 text-sm cursor-grab flex items-center"
          onDragStart={(e) => handleNoteDragStart(e)}
        >
          <RemixIcon icon={riFilePaperLine} className="mr-1" />
          <span>Note</span>
        </div>
      </section>
    </div>
  )
}

export default NodeLibraryPane
