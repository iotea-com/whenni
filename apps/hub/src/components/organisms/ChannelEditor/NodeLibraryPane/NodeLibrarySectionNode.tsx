import { ChannelNode, defaultNodes } from '@iotea/libs/engine/nodes/v1'
import { DragEvent, FC } from 'react'
import styles from './NodeLibrarySectionNode.module.scss'
import { RemixIcon, riDraggable } from '@mwarnerdotme/react-remixicon'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'

type Props = {
  label: string
  defaultNode: ChannelNode
}

const NodeLibrarySectionNode: FC<Props> = ({ label, defaultNode }) => {
  const handleDragStart = (e: DragEvent<HTMLDivElement>, nodeName: string) => {
    const defaultNodeValue = defaultNodes.get(nodeName)
    if (!defaultNodeValue) {
      addToast({
        title: 'No default node confiuration found',
        body: 'If this is a custom node, please add a JSON configuration to the nodes library.',
        level: 'error',
        ttl: -1,
      })
      return
    }
    e.dataTransfer.setData('application/json', JSON.stringify(defaultNodeValue))
  }

  return (
    <div className={styles.sectionNode}>
      <div
        draggable
        onDragStart={(e) => {
          handleDragStart(e, label)
        }}
      >
        <RemixIcon
          className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 dark:text-gray-600"
          icon={riDraggable}
        />
        <h3 className={styles.name}>{defaultNode?.metadata?.name ?? label}</h3>
        {defaultNode?.metadata?.description && (
          <p className={styles.description}>{defaultNode.metadata.description}</p>
        )}
      </div>
    </div>
  )
}

export default NodeLibrarySectionNode
