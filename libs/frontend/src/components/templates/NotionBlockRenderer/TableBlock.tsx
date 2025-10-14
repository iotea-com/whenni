import {
  BlockObjectResponse,
  TableBlockObjectResponse,
} from '@notionhq/client/build/src/api-endpoints'
import RichTextBlock from './RichTextItemBlock'

const TableBlock = ({
  block,
}: {
  block: TableBlockObjectResponse & { children?: BlockObjectResponse[] }
}) => {
  return (
    <div className="overflow-x-auto" key={block.id}>
      <table className="min-w-full border-collapse border border-gray-200">
        <tbody>
          {block.children?.map((row: any, rowIndex: number) => {
            if (row.type !== 'table_row') return null
            return (
              <tr
                key={row.id}
                className={
                  rowIndex === 0 && block.table.has_column_header
                    ? 'bg-gray-50 font-semibold'
                    : 'bg-white'
                }
              >
                {row.table_row.cells.map((cell: any, cellIndex: number) => (
                  <td key={`${row.id}-${cellIndex}`} className="border border-gray-200 px-4 py-2">
                    <RichTextBlock richText={cell} />
                  </td>
                ))}
              </tr>
            )
          })}
        </tbody>
      </table>
    </div>
  )
}

export default TableBlock
