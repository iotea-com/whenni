import { RichTextItemResponse } from '@notionhq/client/build/src/api-endpoints'

const RichTextBlock = ({
  richText,
}: {
  richText: RichTextItemResponse[] | RichTextItemResponse
}) => {
  // Ensure richText is always an array of text objects
  const textArray = (Array.isArray(richText) ? richText : [richText]).filter(
    (v) => v.type === 'text',
  )

  return textArray.map(({ text, annotations }, index) => {
    const { content, link } = text

    // Create className based on annotations
    const className = [
      annotations?.bold ? 'font-bold' : '',
      annotations?.italic ? 'italic' : '',
      annotations?.strikethrough ? 'line-through' : '',
      annotations?.underline ? 'underline' : '',
      annotations?.code ? 'font-mono bg-gray-100 px-1 rounded' : '',
    ]
      .filter(Boolean)
      .join(' ')

    if (link)
      return (
        <a key={index} href={link.url} className={className}>
          {content}
        </a>
      )

    return (
      <span key={index} className={className}>
        {content}
      </span>
    )
  })
}

export default RichTextBlock
