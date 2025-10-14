'use client'

import { FC, useMemo, useRef, useState } from 'react'
import { RemixIcon, IconDefinition, riMore2Fill } from '@mwarnerdotme/react-remixicon'
import { useClickOutside } from '@react-hooks-library/core'
import styles from './index.module.scss'
import { IconSize } from '@mwarnerdotme/react-remixicon/dist/components/RemixIcon'

type ModifierKey = 'ctrl' | 'shift' | 'alt' | 'meta' | 'ctrlOrMeta'
// type ModifierSymbol = '^' | '⌘' | '⇧' | '⌥'

type ContextMenuButton = {
  icon?: IconDefinition
  title: string
  description?: string
  hotkey?: {
    modifiers: ModifierKey[]
    key: string
  }
  onClick: () => void
}

type ProcessedContextMenuButton = Omit<ContextMenuButton, 'hotkey'> & {
  hotkey?: {
    modifiers: string[]
    key: string
  }
}

type Props = {
  contextMenuButtons: ContextMenuButton[]
  iconSize?: IconSize
}

const ContextMenu: FC<Props> = ({ contextMenuButtons, iconSize }) => {
  const [open, setOpen] = useState<boolean>(false)

  const dropdownRef = useRef(null)
  useClickOutside(dropdownRef, () => {
    setOpen(false)
  })

  const toggleOpen = () => {
    setOpen((current) => !current)
  }

  // Determine if we should use meta (macOS) or ctrl (other OS)
  const isCtrlOrMeta = useMemo(() => {
    const isMac =
      // @ts-ignore - userAgentData is not yet in TypeScript's lib.dom
      navigator?.userAgentData?.platform?.toLowerCase().includes('mac') ??
      navigator.platform.toLowerCase().includes('mac')
    return isMac ? 'meta' : 'ctrl'
  }, [])

  // Apply ctrl or meta to the hotkey modifier list
  const processedButtons = useMemo<ProcessedContextMenuButton[]>(() => {
    return contextMenuButtons.map((button) => {
      if (!button || !button.hotkey) return button

      const updatedButton: ProcessedContextMenuButton = { ...button }
      if (updatedButton.hotkey) {
        updatedButton.hotkey = {
          ...updatedButton.hotkey,
          modifiers: updatedButton.hotkey.modifiers.map((modifier) => {
            if (modifier === 'ctrl') return '^'
            if (modifier === 'meta') return '⌘'
            if (modifier === 'shift') return '⇧'
            if (modifier === 'alt') return '⌥'
            if (modifier === 'ctrlOrMeta') return isCtrlOrMeta === 'meta' ? '⌘' : '^'
            return modifier
          }),
        }
      }

      return updatedButton
    })
  }, [contextMenuButtons, isCtrlOrMeta])

  return (
    <div className="relative flex justify-end">
      <RemixIcon
        className="cursor-pointer"
        icon={riMore2Fill}
        size={iconSize}
        onClick={toggleOpen}
      />
      {open && (
        <div ref={dropdownRef} className={styles.contextMenu}>
          {processedButtons.map((contextMenuButton) => (
            <div
              key={contextMenuButton.title}
              className={styles.contextMenuButton}
              onClick={() => {
                contextMenuButton.onClick()
                setOpen(false)
              }}
            >
              {contextMenuButton.hotkey && (
                <div className={styles.hotkey}>
                  <p>
                    {contextMenuButton.hotkey.modifiers.join('') + contextMenuButton.hotkey.key}
                  </p>
                </div>
              )}
              <p className={styles.title}>
                {contextMenuButton.icon && (
                  <RemixIcon className="mr-1" icon={contextMenuButton.icon} />
                )}
                {contextMenuButton.title}
              </p>
              {contextMenuButton.description && (
                <p className={styles.description}>
                  <small>{contextMenuButton.description}</small>
                </p>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default ContextMenu
