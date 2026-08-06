'use client'

import { FC, PropsWithChildren, useMemo } from 'react'
import { Table, flexRender } from '@tanstack/react-table'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import {
  RemixIcon,
  riArrowLeftDoubleLine,
  riArrowLeftSLine,
  riArrowRightDoubleLine,
  riArrowRightSLine,
  riLoader3Line,
  riRefreshLine,
} from '@mwarnerdotme/react-remixicon'
import { useInterval } from '@react-hooks-library/core'

type TableColumnMeta = {
  align?: 'left' | 'center' | 'right'
}

type Props = {
  table: Table<any>
  page?: number
  totalPages?: number
  autoRefresh?: number // must be greater than or equal to 1000 (1 second)
  isLoading?: boolean
  title: string
  description?: string
  handlePageChange?: (newPage: number) => void
  headerContent?: React.ReactNode
}

const ListTable: FC<PropsWithChildren<Props>> = ({
  children,
  table,
  page,
  totalPages,
  autoRefresh,
  isLoading,
  title,
  description,
  handlePageChange,
  headerContent,
}) => {
  const startPage = useMemo(() => {
    if (!page) return 1
    return Math.max(1, page - 2)
  }, [page])

  const endPage = useMemo(() => {
    if (!page || !totalPages) return 1
    return Math.min(totalPages, page + 2)
  }, [page, totalPages])

  const pageNumbers = useMemo(() => {
    let pn: number[] = []
    for (let i = startPage; i <= endPage; i++) {
      pn.push(i)
    }

    return pn
  }, [endPage, startPage])

  useInterval(
    () => {
      if (autoRefresh && autoRefresh >= 1e3 && handlePageChange) handlePageChange(page ?? 1)
    },
    autoRefresh ?? 1e3,
    {
      paused: autoRefresh && autoRefresh >= 1e3 ? false : true,
    },
  )

  return (
    <>
      <div className="mb-2">
        <h2 className="text-lg font-extrabold">{title}</h2>
        {description && <p className="text-xs text-gray-500">{description}</p>}
      </div>
      {headerContent && <div className="mb-2">{headerContent}</div>}
      <div className="bg-gray-100 dark:bg-gray-800 p-[2px] rounded-md">
        <table className="w-full">
          <thead>
            {table.getHeaderGroups().map((headerGroup) => (
              <tr key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  const meta = header.column.columnDef.meta as TableColumnMeta | undefined

                  const alignClass = (() => {
                    switch (meta?.align) {
                      case 'center':
                        return 'text-center'
                      case 'right':
                        return 'text-right'
                      default:
                        return 'text-left'
                    }
                  })()

                  return (
                    <th
                      key={header.id}
                      colSpan={header.colSpan}
                      className={`py-2 px-6 font-normal text-sm text-gray-500 bg-white dark:bg-gray-900 first:rounded-tl last:rounded-tr ${alignClass}`}
                      style={{
                        width: header.column.getIndex() === 0 ? 'auto' : header.column.getSize(),
                      }}
                    >
                      {!header.isPlaceholder &&
                        flexRender(header.column.columnDef.header, header.getContext())}
                    </th>
                  )
                })}
              </tr>
            ))}
          </thead>
          <tbody>
            {table.getRowModel().rows.map((row, rowIndex) => (
              <tr key={row.id} className="group">
                {row.getVisibleCells().map((cell) => {
                  const meta = cell.column.columnDef.meta as TableColumnMeta | undefined

                  const alignClass = (() => {
                    switch (meta?.align) {
                      case 'center':
                        return 'text-center'
                      case 'right':
                        return 'text-right'
                      default:
                        return 'text-left'
                    }
                  })()

                  return (
                    <td
                      key={cell.id}
                      className={`h-12 px-0 first:pl-0 last:pr-0 ${rowIndex === 0 ? 'pt-[2px]' : ''} ${alignClass}`}
                      style={{
                        width: cell.column.getIndex() === 0 ? 'auto' : cell.column.getSize(),
                      }}
                    >
                      <div className="h-full w-full flex items-center transition py-1 px-6 bg-white dark:bg-gray-900 text-gray-700 hover:text-gray-800 dark:text-gray-200 dark:hover:text-gray-100 group-hover:bg-gray-50 dark:group-hover:bg-gray-800">
                        {flexRender(cell.column.columnDef.cell, cell.getContext())}
                      </div>
                    </td>
                  )
                })}
              </tr>
            ))}
            {/* Add empty rows if there are less than 10 rows and page > 1 */}
            {page &&
              page > 1 &&
              table.getRowModel().rows.length < 10 &&
              Array(10 - table.getRowModel().rows.length)
                .fill(0)
                .map((_, index) => (
                  <tr key={`empty-${index}`} className="group">
                    {table.getHeaderGroups()[0].headers.map((header) => (
                      <td
                        key={`empty-${index}-${header.id}`}
                        className="h-12 px-0 first:pl-0 last:pr-0"
                        style={{
                          width: header.column.getIndex() === 0 ? 'auto' : header.column.getSize(),
                        }}
                      >
                        <div className="h-full w-full flex items-center transition py-1 px-6 bg-gray-100 group-hover:bg-gray-100">
                          &nbsp;
                        </div>
                      </td>
                    ))}
                  </tr>
                ))}
          </tbody>
        </table>
        <div className="mt-px bg-white dark:bg-gray-900 rounded-b px-6 py-2 [&_small]:text-gray-500">
          {children}
        </div>
      </div>
      {page && handlePageChange && (
        <div className="flex gap-2 mt-2">
          {Number(totalPages) > 1 && (
            <div className="flex gap-2">
              <Button onClick={() => handlePageChange(1)} disabled={page <= 1} variant="underline">
                <RemixIcon icon={riArrowLeftDoubleLine} />
              </Button>
              <Button
                onClick={() => handlePageChange(page - 1)}
                disabled={page <= 1}
                variant="underline"
              >
                <RemixIcon icon={riArrowLeftSLine} />
              </Button>
              {pageNumbers.map((p) => (
                <Button
                  key={p}
                  onClick={() => handlePageChange(p)}
                  disabled={p === page}
                  variant="underline"
                  className={p === page ? 'text-gray-900!' : ''}
                >
                  {p}
                </Button>
              ))}
              <Button
                onClick={() => handlePageChange(page + 1)}
                disabled={page >= Number(totalPages)}
                variant="underline"
              >
                <RemixIcon icon={riArrowRightSLine} />
              </Button>
              <Button
                onClick={() => handlePageChange(Number(totalPages))}
                disabled={page >= Number(totalPages)}
                variant="underline"
              >
                <RemixIcon icon={riArrowRightDoubleLine} />
              </Button>
            </div>
          )}
          <div className="grow" />
          <Button onClick={() => handlePageChange(page)} variant="underline" disabled={isLoading}>
            {isLoading ? (
              <RemixIcon className="animate-spin dark:text-gray-200" icon={riLoader3Line} />
            ) : (
              <RemixIcon className="dark:text-gray-200" icon={riRefreshLine} />
            )}
          </Button>
        </div>
      )}
    </>
  )
}

export default ListTable
