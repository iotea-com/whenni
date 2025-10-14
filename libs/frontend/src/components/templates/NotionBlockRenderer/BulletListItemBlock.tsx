import {
  BlockObjectResponse,
  BulletedListItemBlockObjectResponse,
} from '@notionhq/client/build/src/api-endpoints'
import RichTextBlock from './RichTextItemBlock'

import styles from './BulletListItemBlock.module.scss'

const BulletListItemBlock = ({
  block,
}: {
  block: BulletedListItemBlockObjectResponse & {
    children?: BlockObjectResponse[]
  }
}) => {
  return (
    <li className={styles.bulletListItem} key={block.id}>
      <span>
        {block.bulleted_list_item && (
          <RichTextBlock richText={block.bulleted_list_item.rich_text} />
        )}
      </span>
      {block.bulleted_list_item && block.has_children && block.children && (
        <ul>
          {block.children.map((child) => {
            if (child.type === 'bulleted_list_item')
              return <BulletListItemBlock key={child.id} block={child} />

            return null
          })}
        </ul>
      )}
    </li>
  )
}

export default BulletListItemBlock
