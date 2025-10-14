import { RemixIcon, riLightbulbFill } from '@mwarnerdotme/react-remixicon'
import { calculateHeadingSlug } from '@iotea/libs/frontend/util/calculateBlogPostSlug'
import { BlockObjectResponse } from '@notionhq/client/build/src/api-endpoints'
import ImageBlock from './ImageBlock'
import RichTextBlock from './RichTextItemBlock'

import NumberListItemBlock from './NumberListItemBlock'
import BulletListItemBlock from './BulletListItemBlock'
import CodeBlock from './CodeBlock'
import './codeBlock.css'
import TableBlock from './TableBlock'

const NotionBlockRenderer = async ({
  blocks,
}: {
  blocks: (BlockObjectResponse & {
    children?: BlockObjectResponse[]
  })[]
}) => {
  return (
    <>
      {blocks.map((block) => {
        switch (block.type) {
          case 'heading_2':
            return (
              <h2
                className="text-gray-800"
                id={calculateHeadingSlug(block.heading_2.rich_text[0].plain_text)}
                key={block.id}
              >
                <RichTextBlock richText={block.heading_2.rich_text} />
              </h2>
            )
          case 'heading_3':
            return (
              <h3
                className="text-gray-800"
                id={calculateHeadingSlug(block.heading_3.rich_text[0].plain_text)}
                key={block.id}
              >
                <RichTextBlock richText={block.heading_3.rich_text} />
              </h3>
            )
          case 'paragraph':
            return (
              <p className="text-gray-800" key={block.id}>
                <RichTextBlock richText={block.paragraph.rich_text} />
              </p>
            )
          case 'image':
            return <ImageBlock key={block.id} block={block} />
          case 'bulleted_list_item':
            return <BulletListItemBlock key={block.id} block={block} />
          case 'numbered_list_item':
            return <NumberListItemBlock key={block.id} block={block} />
          case 'to_do':
            return (
              <div key={block.id}>
                <input type="checkbox" checked={block.to_do.checked} readOnly />
                <span>
                  <RichTextBlock richText={block.to_do.rich_text} />
                </span>
              </div>
            )
          case 'quote':
            return (
              <blockquote key={block.id}>
                <RichTextBlock richText={block.quote.rich_text} />
              </blockquote>
            )
          case 'code':
            return (
              <pre key={block.id} className="codeBlock">
                <CodeBlock
                  content={block.code.rich_text[0].plain_text}
                  language={block.code.language}
                />
              </pre>
            )
          case 'callout':
            return (
              <div
                className="bg-yellow-100 border border-yellow-200 px-6 py-4 text-gray-800 flex gap-4 items-start"
                key={block.id}
              >
                <RemixIcon icon={riLightbulbFill} className="text-yellow-500 shrink-0" size="xl" />
                <div>
                  <RichTextBlock richText={block.callout.rich_text} />
                </div>
              </div>
            )
          case 'table':
            return <TableBlock key={block.id} block={block} />
          default:
            return null
        }
      })}
    </>
  )
}

export default NotionBlockRenderer
