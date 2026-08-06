'use client'

import { FC, useCallback, useEffect, useMemo, useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import Callout from '@gruent/libs/frontend/components/molecules/Callout'
import ChannelEditor from '@gruent/hub/components/organisms/ChannelEditor'
import type { ChannelConfig } from '@gruent/libs/engine/channels/index'
import { useQuery } from '@tanstack/react-query'
import gruentClient from '@gruent/hub/lib/gruent'
import {
  RemixIcon,
  riClipboardFill,
  riLoader3Fill,
  riCloseLine,
  riAddCircleLine,
  riDeleteBin7Line,
  riGitForkFill,
  riUserCommunityFill,
  riCircleLine,
  riHistoryLine,
  riFileSearchLine,
  riShare2Line,
  riSave3Line,
  riLayout3Line,
  riCloudLine,
} from '@mwarnerdotme/react-remixicon'
import useChannelEditorStore from '@gruent/hub/stores/channelEditorStore'
import NodeLibraryPane from '@gruent/hub/components/organisms/ChannelEditor/NodeLibraryPane'
import NodeOptionsPane from '@gruent/hub/components/organisms/ChannelEditor/NodeOptionsPane'
import ChannelOptionsPane from '@gruent/hub/components/organisms/ChannelEditor/ChannelOptionsPane'
import NoteOptionsPane from '@gruent/hub/components/organisms/ChannelEditor/NoteOptionsPane'
import Modal from '@gruent/libs/frontend/components/organisms/Modal'
import ChannelExecutionsTable from '@gruent/hub/components/organisms/ChannelExecutionsTable'
import { addToast } from '@gruent/libs/frontend/hooks/useToast'
import useChannelConfig from '@gruent/libs/frontend/hooks/useChannelConfig'
import copyToClipboard from '@gruent/libs/frontend/util/copyToClipboard'
import { openModal } from '@gruent/libs/frontend/hooks/useModal'
import Link from 'next/link'
import dayjs from 'dayjs'
import Loading from '../loading'
import NodeExecutionLogEntry from '@gruent/libs/frontend/components/molecules/NodeExecutionLogEntry'
import { Space } from '@prisma/client'
import { Organization } from '@prisma/client'
import ContextMenu from '@gruent/hub/components/atoms/ContextMenu'
import { handleDeleteChannel } from '@gruent/hub/actions/channels'
import { useRouter } from 'next/navigation'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import CreateModelModal from '@gruent/hub/components/modals/CreateModelModal'
import CreateThingModal from '@gruent/hub/components/modals/CreateThingModal'
import useAuth from '@gruent/hub/hooks/useAuth'

type Props = {
  orgId: string
  spaceId: string
  channelId: string
  organization: Organization | null
  space: Space | null
}

const ChannelByIdEditClientPage: FC<Props> = ({
  orgId,
  spaceId,
  channelId,
  organization,
  space,
}) => {
  const router = useRouter()

  const [previousPublishStatus, setPreviousPublishStatus] = useState<
    'PUBLISHED' | 'INITIALIZING' | 'UNPUBLISHED'
  >()

  const [channelExecutionsPage, setChannelExecutionsPage] = useState(1)
  const [channelExecutionsTotalPages, setChannelExecutionsTotalPages] = useState(1)
  const [channelExecutionsTotalResults, setChannelExecutionsTotalResults] = useState(0)
  const [channelExecutionsStatusFilter, setChannelExecutionsStatusFilter] = useState<string>()
  const [currentExecutionId, setCurrentExecutionId] = useState<string>()
  const [nodeLibraryOpen, setNodeLibraryOpen] = useState(false)
  const { accessToken } = useAuth()

  const currentNode = useChannelEditorStore((state) => state.currentNode)
  const currentEdge = useChannelEditorStore((state) => state.currentEdge)
  const currentNote = useChannelEditorStore((state) => state.currentNote)
  const upsertNode = useChannelEditorStore((state) => state.upsertNode)
  const upsertEdge = useChannelEditorStore((state) => state.upsertEdge)
  const upsertNote = useChannelEditorStore((state) => state.upsertNote)
  const channelName = useChannelEditorStore((state) => state.channelName)
  const setChannelName = useChannelEditorStore((state) => state.setChannelName)
  const resetChannelEditorStore = useChannelEditorStore((state) => state.reset)
  const runtime = useChannelEditorStore((state) => state.runtime)
  const setRuntime = useChannelEditorStore((state) => state.setRuntime)
  const channelStatus = useChannelEditorStore((state) => state.status)
  const setChannelStatus = useChannelEditorStore((state) => state.setStatus)
  const nodes = useChannelEditorStore((state) => state.nodes)
  const edges = useChannelEditorStore((state) => state.edges)
  const notes = useChannelEditorStore((state) => state.notes)
  const removeNode = useChannelEditorStore((state) => state.removeNode)
  const removeEdge = useChannelEditorStore((state) => state.removeEdge)
  const removeNote = useChannelEditorStore((state) => state.removeNote)
  const setCurrentNode = useChannelEditorStore((state) => state.setCurrentNode)
  const setCurrentEdge = useChannelEditorStore((state) => state.setCurrentEdge)
  const setCurrentNote = useChannelEditorStore((state) => state.setCurrentNote)
  const setValidationErrors = useChannelEditorStore((state) => state.setValidationErrors)
  const setModels = useChannelEditorStore((state) => state.setModels)
  const setThings = useChannelEditorStore((state) => state.setThings)

  const channelConfig = useChannelConfig(channelId, channelName, runtime, nodes, edges, notes)

  const {
    status: getChannelConfigStatus,
    data: channel,
    error: getChannelConfigError,
    refetch: refetchChannel,
  } = useQuery({
    queryKey: ['channels', spaceId, channelId],
    queryFn: async () => {
      if (!accessToken) throw new Error('Invalid auth session.')

      const { data: channel, errors } = await gruentClient(accessToken).channels.get(
        spaceId,
        channelId,
      )

      if (errors && errors.length > 0) throw new Error(errors[0])

      if (!channel) throw new Error('Invalid channel.')

      return channel
    },
  })

  const {
    data: channelExecutionList,
    error: getChannelExecutionsError,
    refetch: refetchChannelExecutions,
    isFetching: isChannelExecutionsFetching,
    isRefetching: isChannelExecutionsRefetching,
  } = useQuery({
    queryKey: [
      'channelExecutions',
      spaceId,
      channelId,
      channelExecutionsPage,
      channelExecutionsStatusFilter,
    ],
    queryFn: async () => {
      if (!accessToken) throw new Error('Invalid auth session.')

      const {
        data: channelExecutions,
        errors,
        totalPages,
        totalResults,
      } = await gruentClient(accessToken).channels.executions.list(spaceId, channelId, {
        page: channelExecutionsPage,
        statusFilter: channelExecutionsStatusFilter,
      })

      if (errors && errors.length > 0) throw new Error(errors[0])

      if (!channelExecutions) throw new Error('Invalid channel executions.')

      setChannelExecutionsTotalPages(totalPages!)
      setChannelExecutionsTotalResults(totalResults!)

      return channelExecutions
    },
  })

  const {
    data: currentExecution,
    error: getCurrentExecutionError,
    isLoading: isCurrentExecutionLoading,
    refetch: refetchCurrentExecution,
  } = useQuery({
    queryKey: ['channelExecution', spaceId, channelId, currentExecutionId],
    enabled: !!currentExecutionId && !!accessToken,
    queryFn: async () => {
      if (!accessToken) throw new Error('Invalid auth session.')
      if (!currentExecutionId) return undefined

      const { data: execution, errors } = await gruentClient(accessToken).channels.executions.get(
        spaceId,
        currentExecutionId,
      )

      if (errors && errors.length > 0) throw new Error(errors[0])
      if (!execution) throw new Error('Invalid execution.')

      return execution
    },
  })

  // Refetch the current execution when the current execution ID changes
  useEffect(() => {
    if (!currentExecutionId) return
    refetchCurrentExecution()
  }, [currentExecutionId, refetchCurrentExecution])

  const currentExecutionDetailsJsx = useMemo(() => {
    if (isCurrentExecutionLoading) return <Loading />
    if (getCurrentExecutionError)
      return (
        <Callout
          className="max-w-xl"
          title="Could not get the selected channel execution"
          description={`${getCurrentExecutionError.message}`}
          variant="error"
        />
      )
    if (!currentExecution) return

    const { nodeExecutionLogs } = currentExecution

    return nodeExecutionLogs.map((nodeExecutionLog, index) => {
      const { level, timestamp, logAttributes } = nodeExecutionLog
      const { node } = logAttributes

      return (
        <NodeExecutionLogEntry
          key={`${node}-${level}-${timestamp}-${index}`}
          nodeExecutionLog={nodeExecutionLog}
        />
      )
    })

    return <></>
  }, [currentExecution, isCurrentExecutionLoading, getCurrentExecutionError])

  const { error: listThingsError, refetch: refetchThings } = useQuery({
    queryKey: ['things', spaceId],
    queryFn: async () => {
      if (!accessToken) throw new Error('Invalid auth session.')

      const { data: things, errors } = await gruentClient(accessToken).things.list(spaceId, {
        resultsPerPage: 100,
      })

      if (errors && errors.length > 0) throw new Error(errors[0])

      if (things) setThings(things)

      return things
    },
  })

  // Add models and  in the space
  const { error: listModelsError, refetch: refetchModels } = useQuery({
    queryKey: ['models', spaceId],
    queryFn: async () => {
      if (!accessToken) throw new Error('Invalid auth session.')

      const { data: models, errors } = await gruentClient(accessToken).models.list(spaceId)

      if (errors && errors.length > 0) throw new Error(errors[0])

      if (models) setModels(models)

      return models
    },
  })

  useEffect(() => {
    if (listThingsError) {
      addToast({
        title: 'Could not retrieve things',
        body: listThingsError.message,
        level: 'error',
      })
    }

    if (listModelsError) {
      addToast({
        title: 'Could not retrieve models',
        body: listModelsError.message,
        level: 'error',
      })
    }
  }, [listThingsError, listModelsError])

  const onThingModalSubmit = useCallback(async () => {
    // Ensure enough time has passed for the thing to be created
    setTimeout(() => {
      refetchThings()
    }, 1000)
  }, [refetchThings])

  const onModelModalSubmit = useCallback(async () => {
    // Ensure enough time has passed for the model to be created
    setTimeout(() => {
      refetchModels()
    }, 1000)
  }, [refetchModels])

  // Reset the channel editor store on initial load
  useEffect(() => {
    resetChannelEditorStore()
  }, [channelId, resetChannelEditorStore])

  // Update nodes and edges whenever the database is queried
  useEffect(() => {
    if (!channel) return
    const channelConfig = JSON.parse(channel.config as string) as ChannelConfig
    const { nodes, edges, notes, runtime } = channelConfig

    // Set nodes
    if (nodes) {
      for (const node of nodes) {
        upsertNode(node.id, node)
      }
    }

    // Set edges
    if (edges) {
      for (const edge of edges) {
        upsertEdge(edge.id, edge)
      }
    }

    // Set notes
    if (notes) {
      for (const note of notes) {
        upsertNote(note.id, note)
      }
    }

    // Set runtime
    setRuntime(runtime)

    // Set name
    setChannelName(channel.name)
  }, [channel, setChannelName, upsertNode, upsertEdge, upsertNote, setRuntime])

  const handleChannelExecutionsPageChange = (newPage: number) => {
    setChannelExecutionsPage(newPage)
    refetchChannelExecutions()
  }

  const handleChannelExecutionsStatusFilterChange = (status: string | null) => {
    if (status === channelExecutionsStatusFilter || status === null) {
      setChannelExecutionsStatusFilter(undefined)
    } else {
      setChannelExecutionsStatusFilter(status)
    }
    refetchChannelExecutions()
  }

  // When channel status changes, refetch the channel and alert the user
  useEffect(() => {
    refetchChannel()

    if (!previousPublishStatus) return
    switch (channelStatus) {
      case 'UNPUBLISHED':
        if (previousPublishStatus === 'PUBLISHED') {
          addToast({
            title: 'Channel unpublished',
            body: 'This channel has been successfully unpublished. It is no longer running and will not respond until it is published again.',
            level: 'info',
          })
        }

        if (previousPublishStatus === 'INITIALIZING') {
          addToast({
            title: 'Channel publish failed',
            body: "This channel failed to initialize. Check the channel's logs for more information.",
            level: 'error',
          })
        }
        break
      case 'PUBLISHED':
        if (previousPublishStatus === 'INITIALIZING') {
          addToast({
            title: 'Channel published',
            body: 'This channel has been successfully published. It is now running and ready to be used!',
            level: 'info',
          })
        }

        if (previousPublishStatus === 'PUBLISHED') {
          addToast({
            title: 'Could not unpublish channel',
            body: 'Something went wrong while unpublishing the channel. If this problem persists, please contact support.',
            level: 'error',
          })
        }
        break
    }
  }, [channelStatus, refetchChannel, previousPublishStatus])

  // Poll for channel status
  useEffect(() => {
    if (!accessToken) return

    const channelStatusPoller = setInterval(() => {
      gruentClient(accessToken)
        .channels.status(spaceId, channelId)
        .then(({ data: newChannelStatus, errors: channelStatusErrors }) => {
          if (channelStatusErrors) {
            addToast({
              title: 'Could not poll for channel status',
              body: channelStatusErrors[0],
              level: 'error',
              ttl: -1,
            })
          }

          if (newChannelStatus && channelStatus !== newChannelStatus.status) {
            setChannelStatus(newChannelStatus.status)
            setPreviousPublishStatus(channelStatus)
          }
        })
    }, 2e3)

    return () => clearInterval(channelStatusPoller)
  }, [channelId, accessToken, spaceId, channelStatus, setChannelStatus])

  // Handle export button
  const handleCopyChannelConfig = useCallback(async () => {
    const { error } = await copyToClipboard(JSON.stringify(channelConfig, null, 2))

    if (error) {
      addToast({
        title: 'Could not copy to clipboard',
        body: 'This may be due to browser permissions. Please allow copying to clipboard and try again.',
        level: 'error',
      })
    }
  }, [channelConfig])

  // handle save button
  const handleSave = useCallback(async () => {
    if (!accessToken) return

    const { errors } = await gruentClient(accessToken).channels.updateConfig(
      spaceId,
      channelConfig.id,
      channelConfig,
    )

    const { data: validationErrors } = await gruentClient(accessToken).channels.validate(
      spaceId,
      channelConfig,
    )

    setValidationErrors(validationErrors)

    if (errors && errors.length > 0) {
      addToast({
        title: 'Could not update the channel',
        body: errors[0],
        level: 'error',
      })

      return
    }

    addToast({
      title: 'Successfully updated the channel',
      body: 'The channel configuration has successfully been saved as a draft.',
      level: 'success',
    })
  }, [channelConfig, spaceId, accessToken, setValidationErrors])

  // handle publish button
  const handlePublish = useCallback(async () => {
    if (!accessToken) return

    const { errors } = await gruentClient(accessToken).channels.publish(spaceId, channelId)

    if (errors && errors.length > 0) {
      addToast({
        title: 'Could not update the channel',
        body: errors[0],
        level: 'error',
      })

      return
    }
  }, [channelId, spaceId, accessToken])

  // handle unpublish button
  const handleUnpublish = useCallback(async () => {
    if (!accessToken) return

    const { errors } = await gruentClient(accessToken).channels.unpublish(spaceId, channelId)

    if (errors && errors.length > 0) {
      addToast({
        title: 'Could not update the channel',
        body: errors[0],
        level: 'error',
      })

      return
    }
  }, [channelId, spaceId, accessToken])

  if (getChannelConfigStatus === 'pending') {
    return (
      <div className="w-full h-full flex justify-center items-center">
        <RemixIcon className="animate-spin" icon={riLoader3Fill} size="3x" />
      </div>
    )
  }

  if (getChannelConfigStatus === 'error') {
    return (
      <div className="w-full h-full flex justify-center items-center">
        <Callout
          className="max-w-xl"
          title="Could not get the channel config"
          description={`${getChannelConfigError.message}`}
          variant="error"
        />
      </div>
    )
  }

  return (
    <>
      <header className="flex border-b border-gray-200 dark:border-gray-800 px-8 py-4 h-14">
        <div className="flex gap-2 items-center w-full text-gray-700 dark:text-gray-300 text-sm">
          <Link
            href={`/organizations/${orgId}`}
            className="flex items-center gap-1 transition text-gray-700 hover:text-gray-800 dark:text-gray-300 dark:hover:text-gray-200"
          >
            <div className="w-5 h-5 flex items-center justify-center bg-gray-100 dark:bg-gray-800 rounded-xs">
              <RemixIcon icon={riUserCommunityFill} size={'sm'} />
            </div>
            {organization?.name}
          </Link>
          <span>/</span>
          <Link
            href={`/organizations/${orgId}/spaces/${spaceId}`}
            className="flex items-center gap-1 transition text-gray-700 hover:text-gray-800 dark:text-gray-300 dark:hover:text-gray-200"
          >
            <div className="w-5 h-5 flex items-center justify-center bg-gray-100 dark:bg-gray-800 rounded-xs">
              <RemixIcon icon={riCircleLine} size={'sm'} />
            </div>
            {space?.name}
          </Link>
          <span>/</span>
          <div className="flex items-center gap-1">
            <div className="w-5 h-5 flex items-center justify-center bg-gray-100 dark:bg-gray-800 rounded-xs">
              <RemixIcon icon={riGitForkFill} size={'sm'} className="rotate-180" />
            </div>
            <Link
              href={`/organizations/${orgId}/spaces/${spaceId}/channels/${channelId}`}
              className="flex items-center gap-1 transition text-gray-700 hover:text-gray-800 dark:text-gray-300 dark:hover:text-gray-200 font-semibold"
            >
              {channel?.name}
            </Link>
          </div>
          <div className="grow" />
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-1">
              <div className="inset-0 flex items-center justify-center text-xs font-medium">
                <AnimatePresence initial={false} mode="wait">
                  {!channelStatus && (
                    <motion.span
                      className="py-1 px-3 rounded-full bg-gray-100 dark:bg-gray-800"
                      initial={{ opacity: 0 }}
                      animate={{ opacity: 1 }}
                      exit={{ opacity: 0 }}
                      transition={{ duration: 0.1 }}
                    >
                      Unknown
                    </motion.span>
                  )}
                  {channelStatus === 'INITIALIZING' && (
                    <motion.div
                      className="flex items-center gap-1"
                      initial={{ opacity: 0 }}
                      animate={{ opacity: 1 }}
                      exit={{ opacity: 0 }}
                      transition={{ duration: 0.1 }}
                    >
                      <RemixIcon icon={riLoader3Fill} size="lg" className="animate-spin" />
                      <span className="py-1 px-3 bg-orange-600 text-white rounded-full">
                        Initializing
                      </span>
                    </motion.div>
                  )}
                  {channelStatus === 'PUBLISHED' && (
                    <motion.span
                      className="py-1 px-3 rounded-full bg-green-600 text-white"
                      initial={{ opacity: 0 }}
                      animate={{ opacity: 1 }}
                      exit={{ opacity: 0 }}
                      transition={{ duration: 0.1 }}
                    >
                      Published
                    </motion.span>
                  )}
                  {channelStatus === 'UNPUBLISHED' && (
                    <motion.span
                      className="py-1 px-3 rounded-full bg-gray-100 dark:bg-gray-800"
                      initial={{ opacity: 0 }}
                      animate={{ opacity: 1 }}
                      exit={{ opacity: 0 }}
                      transition={{ duration: 0.1 }}
                    >
                      Draft
                    </motion.span>
                  )}
                </AnimatePresence>
              </div>
              <button
                className={`
                  relative w-9 h-6 rounded-full transition-colors duration-300 
                  ${channelStatus === 'PUBLISHED' ? 'bg-green-600' : 'bg-gray-300 dark:bg-gray-600'}
                  ${!channelStatus || channelStatus === 'INITIALIZING' ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}
                `}
                onClick={() => {
                  if (!channelStatus || channelStatus === 'INITIALIZING') return
                  channelStatus === 'PUBLISHED' ? handleUnpublish() : handlePublish()
                }}
                disabled={!channelStatus || channelStatus === 'INITIALIZING'}
              >
                <motion.div
                  className="absolute top-1 left-1 w-4 h-4 bg-white rounded-full"
                  animate={{
                    x: channelStatus === 'PUBLISHED' ? 12 : 0,
                  }}
                  transition={{ type: 'spring', stiffness: 500, damping: 30 }}
                />
              </button>
            </div>
            <p className="text-xs text-gray-500">
              Last saved: {dayjs(channel?.updatedAt).format('M/D/YY')} at{' '}
              {dayjs(channel?.updatedAt).format('hh:mm a')}
            </p>
            <Button text="Save" variant="transparent" className="text-xs" onClick={handleSave} />
            <RemixIcon
              icon={riFileSearchLine}
              size="xl"
              className="transition cursor-pointer hover:text-green-700"
              onClick={() => openModal('channelExecutionsTable')}
            />
            <RemixIcon
              icon={riHistoryLine}
              size="xl"
              className="transition cursor-pointer hover:text-green-700"
              onClick={() =>
                addToast({
                  title: 'Coming soon',
                  body: 'Channel versions are coming soon!',
                  level: 'info',
                })
              }
            />
            <ContextMenu
              iconSize="lg"
              contextMenuButtons={[
                {
                  icon: riShare2Line,
                  title: 'Export',
                  hotkey: {
                    modifiers: ['ctrlOrMeta'],
                    key: 'E',
                  },
                  onClick: () => openModal('exportChannel'),
                },
                {
                  icon: riSave3Line,
                  title: 'Save',
                  hotkey: {
                    modifiers: ['ctrlOrMeta'],
                    key: 'S',
                  },
                  onClick: handleSave,
                },
                {
                  icon: riLayout3Line,
                  title: 'Switch modes',
                  hotkey: {
                    modifiers: ['ctrlOrMeta'],
                    key: 'I',
                  },
                  onClick: () => openModal('channelExecutionsTable'),
                },
                {
                  icon: riDeleteBin7Line,
                  title: 'Delete',
                  onClick: async () => {
                    const { error } = await handleDeleteChannel(spaceId, channelId)
                    if (error) {
                      addToast({
                        title: 'Could not delete the channel',
                        body: `The channel could not be deleted: ${error}`,
                        level: 'error',
                      })
                    } else {
                      addToast({
                        title: 'Successfully deleted channel',
                        body: 'Successfully deleted the channel from your space.',
                        level: 'success',
                      })
                      router.push(`/organizations/${orgId}/spaces/${spaceId}`)
                    }
                  },
                },
              ]}
            />
          </div>
        </div>
      </header>
      <div
        id="channelEditor"
        className="flex overflow-y-hidden bg-gray-50 text-gray-700 dark:bg-gray-900 dark:text-gray-200 relative w-[calc(100vw-65px)] h-[calc(100vh-57px)]"
      >
        <Modal id="channelExecutionsTable" showAccept={false}>
          <p>
            <RemixIcon className="mr-1" icon={riCloudLine} />
            {channelStatus === 'PUBLISHED' && (
              <span className="py-1 px-3 text-sm rounded-full bg-green-600 dark:bg-green-300 text-white dark:text-green-900">
                Published
              </span>
            )}
            {channelStatus === 'INITIALIZING' && (
              <span className="py-1 px-3 text-sm bg-orange-600 dark:bg-orange-300 text-white dark:text-orange-900 rounded-full">
                Initializing
              </span>
            )}
            {channelStatus === 'UNPUBLISHED' && (
              <span className="py-1 px-3 text-sm rounded-full bg-gray-100 dark:bg-gray-800">
                Draft
              </span>
            )}
          </p>
          {getChannelExecutionsError && (
            <>
              <Callout
                className="max-w-xl"
                title="Could not get channel executions"
                description={`${getChannelExecutionsError.message}`}
                variant="error"
              />
            </>
          )}
          {!getChannelExecutionsError && (
            <>
              <ChannelExecutionsTable
                orgId={orgId}
                spaceId={spaceId}
                channelExecutionList={channelExecutionList ?? []}
                page={channelExecutionsPage}
                totalPages={channelExecutionsTotalPages}
                totalResults={channelExecutionsTotalResults}
                onPageChange={(newPage) => {
                  handleChannelExecutionsPageChange(newPage)
                }}
                onStatusFilterChange={(status) => {
                  handleChannelExecutionsStatusFilterChange(status)
                }}
                statusFilter={channelExecutionsStatusFilter}
                isLoading={isChannelExecutionsFetching || isChannelExecutionsRefetching}
                autoRefresh={3e3}
                onExecutionClick={(executionId) => {
                  setCurrentExecutionId(executionId)
                  openModal('channelExecutionDetails')
                }}
              />
            </>
          )}
        </Modal>
        {currentExecutionDetailsJsx && (
          <Modal id="channelExecutionDetails" showAccept={false}>
            <div style={{ width: '50vw' }}>
              <div className="mb-4 text-gray-600 dark:text-gray-500">
                <small>
                  <Link
                    className="transition hover:text-gray-800"
                    href={`/organizations/${orgId}/spaces/${spaceId}`}
                  >
                    Space
                  </Link>{' '}
                  &gt;&nbsp;
                  <span
                    className="transition cursor-pointer hover:text-gray-800"
                    onClick={() => {
                      setCurrentExecutionId(undefined)
                      openModal('channelExecutionsTable')
                    }}
                  >
                    Channel
                  </span>
                  &gt;&nbsp;
                  {currentExecutionId}
                </small>
              </div>
              <h2 className="mb-2">Channel Execution Details</h2>
              {currentExecution ? (
                <p>
                  <small>
                    {currentExecution.status} - Started at{' '}
                    {dayjs(currentExecution.startTime).format('MM/DD/YYYY HH:mm:ss.SSS')}, ended at{' '}
                    {dayjs(currentExecution.endTime).format('MM/DD/YYYY HH:mm:ss.SSS')}
                  </small>
                </p>
              ) : (
                <p>
                  <small>Unknown</small>
                </p>
              )}
              {currentExecutionDetailsJsx}
            </div>
          </Modal>
        )}
        <Modal id="exportChannel" showAccept={false}>
          <div
            style={{ maxWidth: 500, maxHeight: 500 }}
            className="overflow-scroll bg-gray-900 text-sm text-gray-200 rounded-sm p-3"
          >
            <pre className="relative">
              <RemixIcon
                icon={riClipboardFill}
                className="absolute top-0 right-0 cursor-pointer hover:text-gray-100"
                onClick={handleCopyChannelConfig}
                size="lg"
              />
              {JSON.stringify(channelConfig, null, 2)}
            </pre>
          </div>
        </Modal>

        <div className="absolute top-5 left-5 flex flex-col gap-1">
          <div className="p-2 cursor-pointer z-30">
            <AnimatePresence initial={false} mode="wait">
              <motion.div
                key={nodeLibraryOpen ? 'close' : 'open'}
                initial={{ opacity: 0, rotate: -90 }}
                animate={{ opacity: 1, rotate: 0 }}
                exit={{ opacity: 0, rotate: 90 }}
                transition={{ duration: 0.1 }}
              >
                {nodeLibraryOpen ? (
                  <RemixIcon
                    className="transition text-gray-500 hover:text-gray-600 dark:text-gray-400 dark:hover:text-gray-300"
                    icon={riCloseLine}
                    size="xl"
                    onClick={() => setNodeLibraryOpen(false)}
                  />
                ) : (
                  <RemixIcon
                    className="transition text-gray-700 dark:text-gray-300 hover:text-primary bg-gray-50 dark:bg-gray-900"
                    icon={riAddCircleLine}
                    size="xl"
                    onClick={() => setNodeLibraryOpen(true)}
                  />
                )}
              </motion.div>
            </AnimatePresence>
          </div>
          <div className="px-2 cursor-pointer z-20">
            <AnimatePresence initial={false} mode="wait">
              <motion.div
                key={
                  currentNode || currentEdge || currentNote
                    ? 'showDeleteSelection'
                    : 'hideDeleteSelection'
                }
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
                transition={{ duration: 0.1 }}
              >
                {(currentNode || currentEdge || currentNote) && (
                  <RemixIcon
                    className="hover:text-red-500"
                    icon={riDeleteBin7Line}
                    size="xl"
                    onClick={() => {
                      if (currentNode) removeNode(currentNode.id)
                      if (currentEdge) removeEdge(currentEdge.id)
                      if (currentNote) removeNote(currentNote.id)
                      setCurrentNode(null)
                      setCurrentEdge(null)
                      setCurrentNote(null)
                    }}
                  />
                )}
              </motion.div>
            </AnimatePresence>
          </div>
        </div>
        <AnimatePresence>
          {nodeLibraryOpen && (
            <motion.div
              className="absolute top-5 left-5 pl-8 pr-2 py-3 w-fit max-w-[300px] bg-gray-50 dark:bg-gray-900 rounded-sm border border-green-700 dark:border-green-300 z-20 overflow-y-scroll"
              style={{ boxShadow: '3px 3px 6px 0 rgba(0, 0, 0, 0.03)', maxHeight: '80%' }}
              initial={{ opacity: 0, scale: 0, transformOrigin: 'top left' }}
              animate={{ opacity: 1, scale: 1, transformOrigin: 'top left' }}
              exit={{ opacity: 0, scale: 0, transformOrigin: 'top left' }}
              transition={{ duration: 0.2, ease: 'easeInOut' }}
            >
              <NodeLibraryPane />
            </motion.div>
          )}
        </AnimatePresence>

        <CreateModelModal spaceId={spaceId} onSubmit={() => onModelModalSubmit()} />
        <CreateThingModal spaceId={spaceId} orgId={orgId} onSubmit={() => onThingModalSubmit()} />

        <ChannelEditor
          hooks={{
            onSave: handleSave,
            onExport: () => openModal('exportChannel'),
            onChangeMode: () => openModal('channelExecutionsTable'),
          }}
        />

        <div
          className="border-l border-gray-200 dark:border-gray-800 flex flex-col overflow-y-scroll"
          style={{ width: '35%', minWidth: 350 }}
        >
          {currentNode && <NodeOptionsPane spaceId={spaceId} orgId={orgId} channelId={channelId} />}
          {currentNote && <NoteOptionsPane />}
          {(currentNode || currentNote) && <hr className="border-gray-200 dark:border-gray-800" />}
          <ChannelOptionsPane spaceId={spaceId} channel={channel} />
        </div>
      </div>
    </>
  )
}

export default ChannelByIdEditClientPage
