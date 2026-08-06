'use client'

import { ChannelNodeType } from '@gruent/libs/engine/nodes/v1'
import { RemixIcon, riArrowRightSLine } from '@mwarnerdotme/react-remixicon'
import { FC, useMemo } from 'react'
import NodeLibrarySectionNode from './NodeLibrarySectionNode'

type Props = {
  nodeType: ChannelNodeType
  open: boolean
  availableNodes: any[]
  onTitleClick: () => void
}

const NodeLibrarySection: FC<Props> = ({ nodeType, open, availableNodes, onTitleClick }) => {
  const sectionNodes = availableNodes.filter(([_nodeName, node]) => {
    if (node.metadata.type === nodeType) return true
    return false
  })

  const sectionTitle = useMemo(() => {
    const titleArray = nodeType.split('')
    titleArray[0] = titleArray[0].toUpperCase()
    return titleArray.join('')
  }, [nodeType])

  if (sectionNodes.length <= 0) return null

  return (
    <section>
      <div className="flex items-center cursor-pointer select-none" onClick={onTitleClick}>
        <RemixIcon
          icon={riArrowRightSLine}
          className={`mr-1 transition transform duration-200 ${open ? 'rotate-90 text-gray-700 dark:text-gray-500' : 'text-gray-600 dark:text-gray-400'}`}
        />
        <h3 className={`transition text-sm font-semibold ${open ? 'text-gray-700' : ''}`}>
          {sectionTitle}
        </h3>
      </div>
      {open && (
        <div>
          {sectionNodes.map(([key, defaultNode]) => {
            // Skip the HTTP Response action node as it is rendered in the HTTP Source node options pane
            if (key === 'HTTP Response-action') return null

            return <NodeLibrarySectionNode key={key} label={key} defaultNode={defaultNode} />
          })}
        </div>
      )}
    </section>
  )
}

export default NodeLibrarySection
