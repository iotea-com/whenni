'use client'

import { openModal } from '@iotea/libs/frontend/hooks/useModal'
import noop from '@iotea/libs/frontend/util/noop'
import Link from 'next/link'
import {
  ButtonHTMLAttributes,
  CSSProperties,
  DetailedHTMLProps,
  FC,
  MouseEvent,
  PropsWithChildren,
  useCallback,
} from 'react'

export type Props = {
  id?: string
  text?: string
  href?: string
  modalId?: string
  type?: DetailedHTMLProps<ButtonHTMLAttributes<HTMLButtonElement>, HTMLButtonElement>['type']
  variant?: 'primary' | 'secondary' | 'transparent' | 'underline' | 'blank'
  className?: string
  style?: CSSProperties
  onClick?: () => void
  disabled?: boolean
}

const Button: FC<PropsWithChildren<Props>> = ({
  id,
  children,
  text,
  href,
  modalId,
  type = 'button',
  variant = 'primary',
  style = {},
  className,
  onClick = noop,
  disabled = false,
}) => {
  const handleClick = useCallback(
    (_e: MouseEvent<HTMLButtonElement>) => {
      // Prevent text selection on button clicks
      if (window.getSelection) window.getSelection()?.removeAllRanges()
      if (modalId) openModal(modalId)
      else onClick()
    },
    [onClick, modalId],
  )

  const blankButton = () => (
    <button
      id={id}
      type={type}
      style={style}
      className={`flex items-center justify-center cursor-pointer rounded-xs border px-3 py-1 font-bold select-none ${className ?? ''}`}
      onClick={handleClick}
      disabled={disabled}
      aria-disabled={disabled}
    >
      {children ?? text}
    </button>
  )

  const primaryButton = () => (
    <button
      id={id}
      type={type}
      style={style}
      className={`flex items-center justify-center cursor-pointer rounded-xs bg-green-600 border border-green-600 px-3 py-1 text-white font-bold disabled:bg-gray-500 disabled:border-gray-500 select-none ${
        className ?? ''
      }`}
      onClick={handleClick}
      disabled={disabled}
      aria-disabled={disabled}
    >
      {children ?? text}
    </button>
  )

  const secondaryButton = () => (
    <button
      id={id}
      type={type}
      style={style}
      className={`flex items-center justify-center cursor-pointer rounded-xs bg-gray-300 border border-gray-300 px-3 py-1 text-gray-800 font-bold select-none ${
        className ?? ''
      }`}
      onClick={handleClick}
      disabled={disabled}
      aria-disabled={disabled}
    >
      {children ?? text}
    </button>
  )

  const underlineButton = () => (
    <span
      id={id}
      style={style}
      className={`cursor-pointer rounded-xs bg-transparent text-gray-700 font-bold underline select-none disabled:text-gray-500 ${
        className ?? ''
      }`}
      onClick={handleClick}
    >
      {children ?? text}
    </span>
  )

  const transparentButton = () => (
    <button
      id={id}
      type={type}
      style={style}
      className={`flex items-center justify-center cursor-pointer rounded-xs bg-transparent px-3 py-1 text-gray-700 border border-gray-700 font-bold hover:border-green-700 hover:text-green-700 transition disabled:border-gray-400 disabled:text-gray-400 dark:text-gray-200 dark:border-gray-200 select-none ${
        className ?? ''
      }`}
      onClick={handleClick}
      disabled={disabled}
      aria-disabled={disabled}
    >
      {children ?? text}
    </button>
  )

  const selectedButton = (() => {
    switch (variant) {
      case 'secondary':
        return secondaryButton()
      case 'underline':
        return underlineButton()
      case 'transparent':
        return transparentButton()
      case 'blank':
        return blankButton()
      case 'primary':
      default:
        return primaryButton()
    }
  })()

  if (href) {
    return (
      <Link href={href} className="cursor-pointer">
        {selectedButton}
      </Link>
    )
  }

  return selectedButton
}

export default Button
