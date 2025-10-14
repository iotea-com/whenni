'use client'

import { useCallback, useEffect, useState } from 'react'
import isBrowser from '../util/isBrowser'

const useModal = (id: string) => {
  // init state
  const checkHash = useCallback(() => {
    if (isBrowser && window.location.hash === `#${id}`) return true
    return false
  }, [id])

  const [isModalOpen, setIsModalOpen] = useState<boolean>(checkHash())

  // check location hash to see if modal is active
  useEffect(() => {
    const ch = () => {
      const hashMatches = checkHash()
      if (hashMatches) setIsModalOpen(true)
      else setIsModalOpen(false)
    }

    if (isBrowser) window.addEventListener('hashchange', ch)
    // return () => window.removeEventListener('hashchange', ch)
  }, [checkHash])

  return isModalOpen
}

/**
 * This function closes a modal by removing the hash from the URL.
 * @returns The `closeModal` function returns `undefined`.
 */
export const closeModal = () => {
  if (!isBrowser) return
  window.location.hash = ' '
  history.replaceState('', document.title, window.location.pathname + window.location.search)
}

/**
 * This function opens a modal by setting the window location hash to a specified ID.
 * @param {string} id - The `id` parameter is a string that represents the unique identifier of the
 * modal that needs to be opened. The function sets the hash of the current URL to the provided `id`,
 * which can be used to show the corresponding modal on the page.
 * @returns The function `openModal` returns nothing (`undefined`). It only sets the
 * `window.location.hash` to the provided `id` if the code is running in a browser environment
 * (`isBrowser` is true).
 */
export const openModal = (id: string) => {
  if (!isBrowser) return
  window.location.hash = id
}

export default useModal
