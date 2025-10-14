'use client'

import { ChangeEvent, FC, useMemo } from 'react'
import noop from '@iotea/libs/frontend/util/noop'

type Props = {
  name: string
  variant?: 'primary' | 'secondary' | 'outline'
  label: string
  hideLabel?: boolean
  defaultValue?: string
  value?: string
  placeholder?: string
  className?: string
  backgroundColor?: string
  width?: 'fit' | 'full'
  height?: 'fit' | 'full'
  disabled?: boolean
  onChange?: (e: ChangeEvent<HTMLTextAreaElement>) => void
  muted?: boolean
}

const FormFieldTextArea: FC<Props> = ({
  name,
  variant = 'outline',
  label,
  hideLabel = false,
  defaultValue,
  value,
  placeholder,
  className = '',
  backgroundColor = 'bg-gray-50',
  width = 'full',
  height = 'full',
  disabled = false,
  onChange = noop,
}) => {
  const calculatedInputClasses = useMemo(() => {
    const inputFieldWidth = (() => {
      switch (width) {
        case 'fit':
          return 'fit-content'
        case 'full':
          return 'w-full'
      }
    })()

    const inputFieldHeight = (() => {
      switch (height) {
        case 'fit':
          return 'fit-content'
        case 'full':
          return 'h-full'
      }
    })()

    const sharedClasses = `border-[1px] rounded px-6 py-3 text-sm outline-none focus:border-green-700 ${inputFieldWidth} ${inputFieldHeight}`

    const variantClasses = (() => {
      switch (variant) {
        case 'outline':
          return `bg-transparent text-gray-100 border-gray-100`
        case 'secondary':
          return ``
        case 'primary':
        default:
          return `border-gray-300 text-gray-800 ${backgroundColor}`
      }
    })()

    return `${sharedClasses} ${variantClasses}`
  }, [variant, width, height, backgroundColor])

  const calculatedLabelClasses = useMemo(() => {
    let sharedClasses = `text-sm px-2 ml-[15px] -top-3 left-0 absolute rounded border-2 border-transparent ${
      hideLabel ? 'hidden' : ''
    }`

    let variantClasses = (() => {
      switch (variant) {
        case 'outline':
          return `bg-black text-gray-400`
        case 'secondary':
          return ``
        case 'primary':
        default:
          return `text-gray-500 ${backgroundColor}`
      }
    })()

    return `${sharedClasses} ${variantClasses}`
  }, [hideLabel, variant, backgroundColor])

  return (
    <div className={`relative my-3 form-field-text-container ${className}`}>
      <textarea
        id={name}
        name={name}
        placeholder={placeholder}
        defaultValue={defaultValue}
        value={value}
        onChange={(e) => onChange(e)}
        className={calculatedInputClasses}
        disabled={disabled}
      />
      <label htmlFor={name} className={calculatedLabelClasses}>
        {label}
      </label>
    </div>
  )
}

export default FormFieldTextArea
