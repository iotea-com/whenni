'use client'

import { FC, PropsWithChildren, useMemo } from 'react'
import { Table, flexRender } from '@tanstack/react-table'
import styles from './index.module.scss'
import Button from '@iotea/libs/frontend/components/atoms/Button'
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
        <table className={styles.table}>
          <thead>
            {table.getHeaderGroups().map((headerGroup) => (
              <tr key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  const meta = header.column.columnDef.meta as TableColumnMeta | undefined

                  const alignClass = (() => {
                    switch (meta?.align) {
                      case 'center':
                        return 'cellAlignCenter'
                      case 'right':
                        return 'cellAlignRight'
                      default:
                        return 'cellAlignLeft'
                    }
                  })()

                  return (
                    <th
                      key={header.id}
                      colSpan={header.colSpan}
                      className={alignClass}
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
            {table.getRowModel().rows.map((row) => (
              <tr key={row.id}>
                {row.getVisibleCells().map((cell) => {
                  const meta = cell.column.columnDef.meta as TableColumnMeta | undefined

                  const alignClass = (() => {
                    switch (meta?.align) {
                      case 'center':
                        return 'cellAlignCenter'
                      case 'right':
                        return 'cellAlignRight'
                      default:
                        return 'cellAlignLeft'
                    }
                  })()

                  return (
                    <td
                      key={cell.id}
                      className={`${styles.cell} ${styles.solidRow} ${alignClass}`}
                      style={{
                        width: cell.column.getIndex() === 0 ? 'auto' : cell.column.getSize(),
                      }}
                    >
                      <div className={styles.cellContent}>
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
                  <tr key={`empty-${index}`} className={styles.emptyRow}>
                    {table.getHeaderGroups()[0].headers.map((header) => (
                      <td
                        key={`empty-${index}-${header.id}`}
                        className={`${styles.cell} ${styles.solidRow}`}
                        style={{
                          width: header.column.getIndex() === 0 ? 'auto' : header.column.getSize(),
                        }}
                      >
                        <div className={styles.cellContent}>&nbsp;</div>
                      </td>
                    ))}
                  </tr>
                ))}
          </tbody>
        </table>
        <div className={styles.stats}>{children}</div>
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
                  className={p === page ? '!text-gray-900' : ''}
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
