'use client'

import useModal, { closeModal } from '@iotea/libs/frontend/hooks/useModal'
import { CSSProperties, FC, PropsWithChildren, useRef } from 'react'

import styles from './index.module.scss'
import noop from '@iotea/libs/frontend/util/noop'
import Button from '../../atoms/Button'
import { RemixIcon, riCloseLine } from '@mwarnerdotme/react-remixicon'
import { useClickOutside } from '@react-hooks-library/core'
import { AnimatePresence, motion } from 'framer-motion'

type Props = {
  id: string
  title?: string
  onAccept?: () => void
  onClose?: () => void
  showClose?: boolean
  showAccept?: boolean
  hideForceClose?: boolean
  acceptText?: string
  style?: CSSProperties
  className?: string
}

const Modal: FC<PropsWithChildren<Props>> = ({
  id,
  children,
  onAccept = noop,
  onClose = closeModal,
  showClose = true,
  showAccept = true,
  hideForceClose = false,
  acceptText = 'Accept',
  style,
  className,
}) => {
  const showModal = useModal(id)

  const modalRef = useRef<HTMLDivElement>(null)
  useClickOutside(modalRef, (e) => {
    // block modal from closing if a toast notification (or element within a toast) is clicked
    const target = e.target as HTMLElement
    if (target && target.closest('#toastNotificationContainer')) return

    // close this modal if it is active
    if (showModal) onClose()
  })

  return (
    <AnimatePresence>
      {showModal && (
        <motion.article
          transition={{ ease: 'easeIn', duration: 0.05 }}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          className={styles.modalWrapper}
        >
          <motion.div
            transition={{ ease: 'easeIn', duration: 0.05 }}
            initial={{ x: 500 }}
            animate={{ x: 0 }}
            exit={{ x: 500 }}
            ref={modalRef}
            style={style}
            className={`${styles.modal} ${className}`}
          >
            {!hideForceClose && (
              <RemixIcon className={styles.closeButton} icon={riCloseLine} onClick={onClose} />
            )}
            <main>{children}</main>
            <div className="grow"></div>
            {(showClose || showAccept) && (
              <footer>
                <hr className="grow border-gray-200 dark:border-gray-800 my-6" />
                <div className="flex gap-4">
                  {showClose && (
                    <Button className="grow" text="Close" onClick={onClose} variant="transparent" />
                  )}
                  {showAccept && (
                    <Button
                      id={id + 'Submit'}
                      className="grow"
                      text={acceptText}
                      onClick={onAccept}
                    />
                  )}
                </div>
              </footer>
            )}
          </motion.div>
        </motion.article>
      )}
    </AnimatePresence>
  )
}

export default Modal
