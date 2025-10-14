'use client'

import { Dispatch, FC, SetStateAction, useCallback, useMemo, useRef, useState } from 'react'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import { Thing, Model, Channel, Tag, AppliedTag } from '@prisma/client'
import { handleApplyTag, handleCreateTag, handleRemoveTag } from '@iotea/hub/actions/tags'
import { RemixIcon, riAddLine, riCloseLine, riLoader2Line } from '@mwarnerdotme/react-remixicon'
import { useClickOutside } from '@react-hooks-library/core'

type Props = {
  spaceId: string
  onClose?: () => void
  onSubmit?: () => void
  subject?: (Thing | Model | Channel) & { tags?: (AppliedTag & { tag: Tag })[] }
  setSubject: Dispatch<SetStateAction<Thing | Model | Channel | undefined>>
  tags: Tag[]
}

const TagMenu: FC<Props> = ({ spaceId, onClose, onSubmit, subject, setSubject, tags }) => {
  const [search, setSearch] = useState('')
  const [isCreateTagLoading, setIsCreateTagLoading] = useState(false)

  const handleClose = useCallback(() => {
    setSearch('')
    setSubject(undefined)

    if (onClose) onClose()
  }, [onClose, setSubject])

  const menuRef = useRef(null)
  useClickOutside(menuRef, () => {
    handleClose()
  })

  const handleClickRemoveTag = useCallback(
    async (tagId: string) => {
      if (!subject) return

      const { error: removeTagError } = await handleRemoveTag(spaceId, tagId, subject.id)
      if (removeTagError) {
        addToast({
          title: 'Could not remove tag',
          body: `${removeTagError}`,
          level: 'error',
        })
        return
      }
    },
    [spaceId, subject],
  )

  const handleClickAvailableTag = useCallback(
    async (spaceId: string, tagId: string) => {
      if (!subject) return

      const { error: applyTagError } = await handleApplyTag(spaceId, tagId, subject.id)
      if (applyTagError) {
        addToast({
          title: 'Could not apply tag',
          body: `${applyTagError}`,
          level: 'error',
        })
        return
      }

      handleClose()
    },
    [subject, handleClose],
  )

  const handleNewTag = useCallback(async () => {
    if (!subject) return

    setIsCreateTagLoading(true)

    if (!search) {
      addToast({
        title: 'Could not create tag',
        body: 'A tag name is required',
        level: 'warning',
      })
      setIsCreateTagLoading(false)
      return
    }

    const { data, error: createTagError } = await handleCreateTag(spaceId, search)
    if (createTagError) {
      addToast({
        title: 'Could not create tag',
        body: `${createTagError}`,
        level: 'error',
      })
      setIsCreateTagLoading(false)
      return
    }

    if (!data) {
      addToast({
        title: 'Could not create tag',
        body: 'An unknown error occurred while creating the tag. No tag was returned from the server. Please contact support.',
        level: 'error',
      })
      setIsCreateTagLoading(false)
      return
    }

    const { tag } = data
    const { error: applyTagError } = await handleApplyTag(spaceId, tag.id, subject.id)

    if (applyTagError) {
      addToast({
        title: 'Could not add tag',
        body: `${applyTagError}`,
        level: 'error',
      })
      setIsCreateTagLoading(false)
      return
    }

    setIsCreateTagLoading(false)
    handleClose()
    if (onSubmit) onSubmit()
  }, [search, spaceId, handleClose, onSubmit, subject])

  const usedTags = useMemo(() => {
    return tags
      .map((tag) => {
        if (!subject || !subject.tags) return null
        if (subject.tags.find((appliedTag) => appliedTag.tagId === tag.id)) return tag
        return null
      })
      .filter((tag) => tag !== null)
  }, [tags, subject])

  const availableTags = useMemo(() => {
    if (!search) return []

    const filteredTags = tags.filter((tag) => {
      // Filter out tags that are already applied to the subject
      if (usedTags.find((usedTag) => usedTag.id === tag.id)) return false

      // Filter out tags that don't match the search query
      if (!tag.name.toLowerCase().includes(search.toLowerCase())) return false

      return true
    })

    return filteredTags.slice(0, 5)
  }, [tags, usedTags, search])

  const shouldShowCreateTag = useMemo(() => {
    return search.length > 0 && !tags.find((tag) => tag.name.toLowerCase() === search.toLowerCase())
  }, [search, tags])

  return (
    <>
      {subject && (
        <div
          ref={menuRef}
          className="absolute w-[300px] top-0 left-0 w-fit p-3 border border-gray-300 dark:border-gray-700 bg-gray-200 dark:bg-gray-900 shadow-md rounded-md z-10"
        >
          <div className="absolute top-0 right-1">
            <RemixIcon
              className="cursor-pointer"
              icon={riCloseLine}
              size="sm"
              onClick={() => setSubject(undefined)}
            />
          </div>
          <div className="flex flex-wrap gap-2 items-center">
            {usedTags.map((tag) => (
              <div
                key={tag.id}
                className="flex items-center gap-[2px] pl-2 pr-1 py-1 bg-gray-300 dark:bg-gray-700 text-gray-800 dark:text-gray-200 text-sm rounded-full"
              >
                {tag.name}
                <button
                  type="button"
                  className="w-4 h-4 flex items-center justify-center transition hover:bg-gray-400 hover:dark:bg-gray-600 rounded-full text-xs hover:bg-gray-500 dark:hover:bg-gray-500"
                  onClick={() => handleClickRemoveTag(tag.id)}
                >
                  <RemixIcon icon={riCloseLine} size="sm" />
                </button>
              </div>
            ))}
            <input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="flex-1 text-sm min-w-[100px] bg-transparent border-none focus:outline-none"
              placeholder="Add tags..."
            />
          </div>
          <div className="flex flex-col gap-2 mt-2">
            {availableTags.map((tag) => (
              <div
                key={tag.id}
                className="px-2 py-1 cursor-pointer text-sm text-gray-800 dark:text-gray-200 transition rounded block border border-transparent hover:bg-green-100 dark:hover:bg-green-800 hover:border-green-200 dark:hover:border-green-700 hover:text-green-700 dark:hover:text-green-300"
                onClick={() => handleClickAvailableTag(spaceId, tag.id)}
              >
                {tag.name}
              </div>
            ))}
            {shouldShowCreateTag && isCreateTagLoading && (
              <div className="px-2 py-1 cursor-pointer text-sm text-gray-800 dark:text-gray-200 transition rounded block">
                <RemixIcon icon={riLoader2Line} className="animate-spin" size="sm" />
              </div>
            )}
            {shouldShowCreateTag && !isCreateTagLoading && (
              <div
                onClick={() => handleNewTag()}
                className="px-2 py-1 cursor-pointer text-sm text-gray-800 dark:text-gray-200 transition rounded block border border-transparent hover:bg-green-100 dark:hover:bg-green-800 hover:border-green-200 dark:hover:border-green-700 hover:text-green-700 dark:hover:text-green-300"
              >
                Create {search} <RemixIcon icon={riAddLine} size="sm" />
              </div>
            )}
          </div>
        </div>
      )}
    </>
  )
}

export default TagMenu
