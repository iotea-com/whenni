import { calculateHeadingSlug } from '@gruent/libs/frontend/util/calculateBlogPostSlug'
import Link from 'next/link'

type Props = {
  tableOfContents: Record<
    string,
    {
      subheadings: Record<
        string,
        {
          type: string
        }
      >
    }
  >
  className?: string
}

const TableOfContents = ({ tableOfContents, className }: Props) => {
  return (
    <aside
      className={`h-fit bg-gray-50 border border-gray-200 rounded p-4 mr-4 col-span-12 md:col-span-4 lg:col-span-3 max-h-[calc(100vh-5rem)] overflow-y-scroll text-xs ${className}`}
    >
      <nav className="space-y-1">
        {Object.entries(tableOfContents).map(([heading, content]) => (
          <div key={heading}>
            <Link
              className="text-left w-full text-gray-600 transition hover:text-green-700 select-none"
              href={`#${calculateHeadingSlug(heading)}`}
            >
              {heading}
            </Link>
            {content.subheadings && (
              <ul className="mt-1 ml-4 space-y-1">
                {Object.entries(content.subheadings).map(([subheading, _]) => (
                  <li key={subheading}>
                    <Link
                      className="text-left w-full text-gray-600 transition hover:text-green-700 select-none"
                      href={`#${calculateHeadingSlug(subheading)}`}
                    >
                      {subheading}
                    </Link>
                  </li>
                ))}
              </ul>
            )}
          </div>
        ))}
      </nav>
    </aside>
  )
}

export default TableOfContents
