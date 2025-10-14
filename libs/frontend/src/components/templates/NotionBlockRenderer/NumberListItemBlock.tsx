import {
  BlockObjectResponse,
  NumberedListItemBlockObjectResponse,
} from '@notionhq/client/build/src/api-endpoints'
import RichTextBlock from './RichTextItemBlock'

import styles from './NumberListItemBlock.module.scss'

const NumberListItemBlock = ({
  block,
}: {
  block: NumberedListItemBlockObjectResponse & {
    children?: BlockObjectResponse[]
  }
}) => {
  return (
    <li className={styles.numberListItem} key={block.id}>
      <span>
        {block.numbered_list_item && (
          <RichTextBlock richText={block.numbered_list_item.rich_text} />
        )}
      </span>
      {block.numbered_list_item && block.has_children && block.children && (
        <ol>
          {block.children.map((child) => {
            if (child.type === 'numbered_list_item')
              return <NumberListItemBlock key={child.id} block={child} />

            return null
          })}
        </ol>
      )}
    </li>
  )
}

export default NumberListItemBlock
