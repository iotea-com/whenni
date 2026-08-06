'use client'

import { FC, useCallback, useState } from 'react'
import {
  ColumnDef,
  createColumnHelper,
  getCoreRowModel,
  getPaginationRowModel,
  useReactTable,
} from '@tanstack/react-table'
import ListTable from '@gruent/hub/components/atoms/ListTable'
import SecretContextMenu from './SecretContextMenu'
import { riInformationLine } from '@mwarnerdotme/react-remixicon'
import { RemixIcon } from '@mwarnerdotme/react-remixicon'
// import IndeterminateCheckbox from "../../atoms/IndeterminateCheckbox"

type Props = {
  secrets: { name: string }[]
  spaceId: string
  // page: number
  totalPages: number
  totalResults: number
}

const SecretsTable: FC<Props> = ({ secrets, spaceId, totalPages, totalResults }) => {
  const [currentPage, setCurrentPage] = useState(1)

  const [rowSelection, setRowSelection] = useState({})

  const columnHelper = createColumnHelper<{ name: string }>()

  const columns: ColumnDef<{ name: string }>[] = [
    // TODO: enable selection when there is bulk action support in the UI for selected rows
    // {
    //   id: 'select',
    //   cell: ({ row }) => (
    //     <IndeterminateCheckbox
    //       {...{
    //         checked: row.getIsSelected(),
    //         disabled: !row.getCanSelect(),
    //         indeterminate: row.getIsSomeSelected(),
    //         onChange: row.getToggleSelectedHandler(),
    //       }}
    //     />
    //   ),
    // },
    columnHelper.accessor('name', {
      header: () => (
        <div className="flex items-center">
          <RemixIcon className="mr-1" icon={riInformationLine} />
          <span>Name</span>
        </div>
      ),
      cell: (info) => info.row.original.name,
      enableSorting: true,
    }),
    columnHelper.display({
      id: 'contextMenu',
      cell: (info) => <SecretContextMenu spaceId={spaceId} secretName={info.row.original.name} />,
      size: 25,
    }),
  ]

  const table = useReactTable<{ name: string }>({
    data: secrets,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getRowId: (original) => original.name,
    enableRowSelection: true,
    // manualPagination: true,
    state: {
      rowSelection,
    },
    onRowSelectionChange: setRowSelection,
  })

  const handlePageChange = useCallback(
    (newPage: number) => {
      if (newPage > totalPages) {
        table.setPageIndex(totalPages - 1)
        setCurrentPage(totalPages)
      } else if (newPage < 1) {
        table.setPageIndex(0)
        setCurrentPage(1)
      } else {
        table.setPageIndex(newPage - 1)
        setCurrentPage(newPage)
      }
    },
    [table, totalPages],
  )

  if (totalResults <= 0) return <p>No secrets found for this space</p>

  return (
    <ListTable
      title="Secrets"
      description="Keys, passwords, tokens - anything that you need to stay private."
      table={table}
      page={currentPage}
      totalPages={totalPages}
      handlePageChange={handlePageChange}
      autoRefresh={9e9}
    >
      {totalResults === 1 ? (
        <small>{totalResults} secret in this space</small>
      ) : (
        <small>{totalResults} secrets in this space</small>
      )}
    </ListTable>
  )
}

export default SecretsTable
