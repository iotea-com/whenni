'use client'

import { FC, useCallback, useRef, useState } from 'react'
import {
  ColumnDef,
  createColumnHelper,
  getCoreRowModel,
  getPaginationRowModel,
  useReactTable,
} from '@tanstack/react-table'
import { AppliedTag, Model, Tag } from '@prisma/client'
import ListTable from '@iotea/hub/components/atoms/ListTable'
import Link from 'next/link'
import ModelContextMenu from './ModelContextMenu'
import { useRouter } from 'next/navigation'
import {
  RemixIcon,
  riAddCircleLine,
  riAddLine,
  riBookShelfLine,
  riEqualizerLine,
} from '@mwarnerdotme/react-remixicon'
import { riInformationLine } from '@mwarnerdotme/react-remixicon'
import TagMenu from '../../modals/TagMenu'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import { useClickOutside } from '@react-hooks-library/core'
import SearchBar from '../../molecules/SearchBar'
// import IndeterminateCheckbox from "../../atoms/IndeterminateCheckbox"

type Props = {
  models: (Model & { tags?: (AppliedTag & { tag: Tag })[] })[]
  tags: Tag[]
  orgId: string
  spaceId: string
  page: number
  totalPages: number
  totalResults: number
  searchFilter?: string
  tagFilter?: string[]
}

const ModelsTable: FC<Props> = ({
  models,
  tags,
  orgId,
  spaceId,
  page,
  totalPages,
  totalResults,
  searchFilter,
  tagFilter,
}) => {
  const [rowSelection, setRowSelection] = useState({})
  const [tagCreateSubject, setTagCreateSubject] = useState<Model>()
  const [showFilterOptions, setShowFilterOptions] = useState(false)
  const router = useRouter()
  const filterOptionsRef = useRef<HTMLDivElement>(null)

  useClickOutside(filterOptionsRef, () => setShowFilterOptions(false))

  const columnHelper = createColumnHelper<Model & { tags?: (AppliedTag & { tag: Tag })[] }>()

  const columns: ColumnDef<Model & { tags?: (AppliedTag & { tag: Tag })[] }>[] = [
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
      cell: (info) => (
        <Link href={`/organizations/${orgId}/spaces/${spaceId}/models/${info.row.id}`}>
          {info.getValue()}
        </Link>
      ),
      enableSorting: true,
    }),
    columnHelper.accessor('tags', {
      header: () => (
        <div className="flex items-center">
          <RemixIcon className="mr-1" icon={riBookShelfLine} />
          <span>Tags</span>
        </div>
      ),
      cell: (info) => {
        const value = info.getValue()
        if (value && value.length > 0) {
          return (
            <div className="relative flex flex-wrap gap-2">
              {value.map((appliedTag, index) => {
                if (index > 3) return null
                if (index > 2)
                  return (
                    <span
                      key="__ellipsis__"
                      className="px-3 transition bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200 rounded-full text-sm"
                    >
                      ...
                    </span>
                  )
                return (
                  <Link
                    href={`/organizations/${orgId}/spaces/${spaceId}/tags/${appliedTag.tagId}`}
                    key={appliedTag.id}
                  >
                    <span className="px-2 py-1 transition bg-gray-100 hover:bg-gray-200 text-gray-800 dark:bg-gray-700 dark:hover:bg-gray-600 dark:text-gray-200 rounded-full text-sm">
                      {appliedTag.tag.name}
                    </span>
                  </Link>
                )
              })}
              {tagCreateSubject === info.row.original && (
                <TagMenu
                  spaceId={spaceId}
                  subject={info.row.original}
                  setSubject={setTagCreateSubject}
                  tags={tags}
                />
              )}
              <Button variant="underline" onClick={() => setTagCreateSubject(info.row.original)}>
                <RemixIcon icon={riAddLine} size="sm" />
              </Button>
            </div>
          )
        }
        return (
          <div className="relative">
            <Button variant="underline" onClick={() => setTagCreateSubject(info.row.original)}>
              <RemixIcon icon={riAddLine} size="sm" />
            </Button>
            {tagCreateSubject === info.row.original && (
              <TagMenu
                spaceId={spaceId}
                subject={info.row.original}
                setSubject={setTagCreateSubject}
                tags={tags}
              />
            )}
          </div>
        )
      },
    }),
    columnHelper.display({
      id: 'contextMenu',
      cell: (info) => <ModelContextMenu spaceId={spaceId} modelId={info.row.original.id} />,
      size: 50,
      meta: {
        align: 'right',
      },
    }),
  ]

  const handleTagFilterChange = useCallback(
    (tagId: string) => {
      const updatedTagFilter = (() => {
        if (!tagFilter) return [tagId]

        if (tagFilter.includes(tagId)) return tagFilter.filter((t) => t !== tagId)
        return [...tagFilter, tagId]
      })()

      const newLocation = new URL(window.location.toString())
      if (updatedTagFilter.length > 0) {
        newLocation.searchParams.set('tagFilter', updatedTagFilter.join(','))
      } else {
        newLocation.searchParams.delete('tagFilter')
      }
      router.replace(newLocation.toString(), { scroll: false })
    },
    [router, tagFilter],
  )

  const clearFilters = useCallback(() => {
    // Clear tag filter
    const newLocation = new URL(window.location.toString())
    newLocation.searchParams.delete('tagFilter')
    router.replace(newLocation.toString(), { scroll: false })
  }, [router])

  const table = useReactTable<Model & { tags?: (AppliedTag & { tag: Tag })[] }>({
    data: models,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getRowId: (original) => original.id,
    enableRowSelection: true,
    manualPagination: true,
    state: {
      rowSelection,
    },
    onRowSelectionChange: setRowSelection,
  })

  const handlePageChange = (newPage: number) => {
    const newLocation = new URL(window.location.toString())
    newLocation.searchParams.set('page', newPage.toString())
    router.replace(newLocation.toString(), { scroll: false })
  }

  return (
    <ListTable
      title="Models"
      description="Define the structure of your data."
      table={table}
      page={page}
      totalPages={totalPages}
      handlePageChange={handlePageChange}
      headerContent={
        <>
          <div className="flex flex-wrap items-end gap-2 mb-1 relative">
            <SearchBar paramKey="q" initialFilter={searchFilter} placeholder="Search things" />
            <div className="grow" />
            <Button
              variant="transparent"
              className="px-2"
              onClick={() => setShowFilterOptions((current) => !current)}
            >
              <RemixIcon icon={riEqualizerLine} size="lg" />
            </Button>
            <Button variant="transparent" className="text-sm" modalId="createModel">
              <RemixIcon className="mr-1" icon={riAddCircleLine} />
              <span>Create</span>
            </Button>
            {showFilterOptions && (
              <div
                ref={filterOptionsRef}
                className="absolute max-w-[500px] z-10 top-10 right-0 bg-gray-100 dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-md p-4 flex flex-col gap-2 overflow-hidden"
              >
                <div className="border-b border-gray-200 dark:border-gray-800 pb-2">
                  <h2 className="text-sm font-medium mb-1">Tags</h2>
                  <div className="flex flex-wrap items-center gap-2">
                    {tags.map((tag) => {
                      if (tagFilter && tagFilter.includes(tag.id))
                        return (
                          <Button
                            key={tag.id}
                            variant="transparent"
                            className="text-sm !border-primary"
                            onClick={() => handleTagFilterChange(tag.id)}
                          >
                            {tag.name}
                          </Button>
                        )

                      return (
                        <Button
                          key={tag.id}
                          variant="transparent"
                          className="text-sm"
                          onClick={() => handleTagFilterChange(tag.id)}
                        >
                          {tag.name}
                        </Button>
                      )
                    })}
                  </div>
                </div>
                <span
                  className="cursor-pointer transition text-xs text-gray-600 hover:text-red-500"
                  onClick={clearFilters}
                >
                  Clear filters
                </span>
              </div>
            )}
          </div>
        </>
      }
    >
      {totalResults === 1 ? (
        <small>{totalResults} model in this space</small>
      ) : (
        <small>{totalResults} models in this space</small>
      )}
    </ListTable>
  )
}

export default ModelsTable
