'use client'

import { ChangeEvent, FC, useMemo } from 'react'
import noop from '@gruent/libs/frontend/util/noop'

type Props = {
  name: string
  variant?: 'primary' | 'secondary' | 'outline-solid'
  inputType?: 'text' | 'password' | 'number' | 'email'
  min?: string
  max?: string
  label: string
  hideLabel?: boolean
  defaultValue?: string
  value?: string
  placeholder?: string
  className?: string
  backgroundColor?: string
  width?: 'fit' | 'full'
  disabled?: boolean
  onChange?: (e: ChangeEvent<HTMLInputElement>) => void
  autocomplete?: 'off' | 'on'
  muted?: boolean
  spellCheck?: boolean
}

const FormFieldText: FC<Props> = ({
  name,
  variant = 'primary',
  inputType = 'text',
  min,
  max,
  label,
  hideLabel = false,
  defaultValue,
  value,
  placeholder,
  className = '',
  backgroundColor = 'bg-gray-50 dark:bg-gray-900',
  width = 'full',
  disabled = false,
  onChange = noop,
  autocomplete = 'off',
  spellCheck = false,
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

    const sharedClasses = `border rounded-xs px-6 py-3 text-sm ${inputFieldWidth}`

    const variantClasses = (() => {
      switch (variant) {
        case 'outline-solid':
          return `bg-transparent text-gray-100 border-gray-100 transition outline-hidden focus:border-green-700`
        case 'secondary':
          return ``
        case 'primary':
        default:
          return `border-gray-300 text-gray-800 dark:text-gray-200 transition outline-hidden focus:border-green-700 ${backgroundColor}`
      }
    })()

    return `${sharedClasses} ${variantClasses}`
  }, [variant, width, backgroundColor])

  const calculatedLabelClasses = useMemo(() => {
    let sharedClasses = `text-sm px-2 ml-[15px] -top-3 left-0 absolute rounded border-2 border-transparent ${
      hideLabel ? 'hidden' : ''
    }`

    let variantClasses = (() => {
      switch (variant) {
        case 'outline-solid':
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
      <input
        id={name}
        name={name}
        type={inputType}
        min={min}
        max={max}
        placeholder={placeholder}
        defaultValue={defaultValue}
        value={value}
        onChange={(e) => onChange(e)}
        className={calculatedInputClasses}
        readOnly={disabled}
        aria-readonly={disabled}
        autoComplete={autocomplete}
        spellCheck={spellCheck}
      />
      <label htmlFor={name} className={calculatedLabelClasses}>
        {label}
      </label>
    </div>
  )
}

export default FormFieldText
