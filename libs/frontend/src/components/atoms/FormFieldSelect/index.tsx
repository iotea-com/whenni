'use client'

import { ChangeEvent, FC, useMemo, useRef, useState, useEffect } from 'react'
import noop from '@iotea/libs/frontend/util/noop'
import {
  RemixIcon,
  riArrowDownSLine,
  riArrowUpSLine,
  riCheckboxBlankCircleLine,
  riCheckboxCircleFill,
  riCloseLine,
} from '@mwarnerdotme/react-remixicon'
import { motion, AnimatePresence } from 'framer-motion'

type Option = {
  value: any
  label: string
  description?: string // only used in cards variant
  icon?: string // only used in cards variant
}

type Props = {
  name: string
  variant?: 'primary' | 'secondary' | 'outline' | 'cards'
  options: Option[]
  optional?: boolean
  label: string
  hideLabel?: boolean
  hideOptionalLabel?: boolean
  className?: string
  width?: 'fit' | 'full'
  multiple?: boolean
  size?: number
  defaultValue?: string | number | readonly string[]
  value?: string | number | readonly string[]
  disabled?: boolean
  backgroundColor?: string
  onChange?: (e: ChangeEvent<HTMLSelectElement>) => void
}

const FormFieldSelect: FC<Props> = ({
  name,
  variant = 'primary',
  options = [],
  optional = false,
  label,
  hideLabel = false,
  hideOptionalLabel = false,
  className = '',
  width = 'full',
  multiple = false,
  defaultValue,
  value,
  disabled = false,
  backgroundColor = 'bg-gray-50 dark:bg-gray-900',
  onChange = noop,
}) => {
  const [isOpen, setIsOpen] = useState(false)
  const [selectedOptions, setSelectedOptions] = useState<Option[]>([])
  const dropdownRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    // First check for controlled value prop
    if (value !== undefined) {
      if (!multiple) {
        const singleOption = options.find((opt) => opt.value === value)
        setSelectedOptions(singleOption ? [singleOption] : [])
        return
      }
      // Handle array of values for multiple select
      const controlledValues = Array.isArray(value) ? value : []
      setSelectedOptions(options.filter((opt) => controlledValues.includes(opt.value)))
      return
    }

    // If no value prop, use defaultValue
    if (defaultValue !== undefined) {
      if (!multiple) {
        const singleOption = options.find((opt) => opt.value === defaultValue)
        setSelectedOptions(singleOption ? [singleOption] : [])
        return
      }
      // Handle array of defaultValues for multiple select
      const defaultValues = Array.isArray(defaultValue) ? defaultValue : []
      setSelectedOptions(options.filter((opt) => defaultValues.includes(opt.value)))
      return
    }
  }, [value, defaultValue, multiple, options])

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  const handleOptionClick = (option: Option) => {
    let newSelectedOptions: Option[]

    if (multiple) {
      if (selectedOptions.some((opt) => opt.value === option.value)) {
        // Remove option if already selected
        newSelectedOptions = selectedOptions.filter((opt) => opt.value !== option.value)
      } else {
        // Add option to selection
        newSelectedOptions = [...selectedOptions, option]
      }
    } else {
      if (optional && selectedOptions.some((opt) => opt.value === option.value)) {
        // Remove option if already selected
        newSelectedOptions = []
      } else {
        // Add option to selection
        newSelectedOptions = [option]
      }
      setIsOpen(false)
    }

    setSelectedOptions(newSelectedOptions)

    // Simulate select onChange event
    const simulatedEvent = {
      target: {
        name,
        value: multiple
          ? newSelectedOptions.map((opt) => opt.value)
          : (newSelectedOptions[0]?.value ?? null),
      },
    } as ChangeEvent<HTMLSelectElement>
    onChange(simulatedEvent)
  }

  const handleClearSelection = () => {
    setSelectedOptions([])
    const simulatedEvent = {
      target: {
        name,
        value: multiple ? [] : '',
      },
    } as ChangeEvent<HTMLSelectElement>
    onChange(simulatedEvent)
  }

  // Handle keyboard navigation
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (disabled) return

    switch (e.key) {
      case 'Enter':
      case 'Space':
        e.preventDefault()
        setIsOpen(!isOpen)
        break
      case 'Tab':
        setIsOpen(false)
        break
      case 'Escape':
        setIsOpen(false)
        break
      case 'ArrowUp':
        e.preventDefault()
        if (!isOpen) setIsOpen(true)
        const nextOptionUp = options[options.indexOf(selectedOptions[0] || options[0]) - 1]
        if (nextOptionUp) setSelectedOptions([nextOptionUp])
        break
      case 'ArrowDown':
        e.preventDefault()
        if (!isOpen) setIsOpen(true)
        const nextOptionDown = options[options.indexOf(selectedOptions[0] || options[0]) + 1]
        if (nextOptionDown) setSelectedOptions([nextOptionDown])
        break
    }
  }

  const calculatedInputClasses = useMemo(() => {
    const inputFieldWidth = (() => {
      switch (width) {
        case 'fit':
          return 'fit-content'
        case 'full':
          return 'w-full'
      }
    })()

    const sharedClasses = `border-[1px] rounded-sm px-6 py-3 h-auto appearance-none cursor-pointer ${inputFieldWidth}`

    const variantClasses = (() => {
      switch (variant) {
        case 'outline':
          return `bg-transparent text-gray-700 border-gray-300 active:border-green-700 focus:border-green-700 outline-none`
        case 'secondary':
          return ``
        case 'primary':
        default:
          return `border-gray-300 text-gray-700 bg-transparent active:border-green-700 focus:border-green-700 outline-none`
      }
    })()

    return `${sharedClasses} ${variantClasses}`
  }, [variant, width])

  const calculatedLabelClasses = useMemo(() => {
    const sharedClasses = `z-10 px-2 ml-[15px] -top-3 left-0 text-sm absolute rounded border-2 border-transparent select-none ${
      hideLabel ? 'hidden' : ''
    }`

    const variantClasses = (() => {
      switch (variant) {
        case 'outline':
          return `text-gray-400 ${backgroundColor}`
        case 'secondary':
          return ``
        case 'primary':
        default:
          return `text-gray-500 ${backgroundColor}`
      }
    })()

    return `${sharedClasses} ${variantClasses}`
  }, [hideLabel, variant, backgroundColor])

  if (variant === 'cards') {
    return (
      <div>
        {/* Add hidden select element for form data */}
        <select
          name={name}
          value={
            multiple ? selectedOptions.map((opt) => opt.value) : (selectedOptions[0]?.value ?? '')
          }
          onChange={onChange}
          multiple={multiple}
          hidden
          aria-hidden="true"
          disabled={disabled}
        >
          {options.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
        <label className={`text-sm text-gray-600 dark:text-gray-400 ${hideLabel ? 'hidden' : ''}`}>
          {label}
          {optional && !hideOptionalLabel ? ' (optional)' : ''}
        </label>
        <div className={`flex flex-wrap gap-2 relative mt-1 mb-3 ${className}`}>
          {options.map((option) => {
            const isSelected = selectedOptions.some((opt) => opt.value === option.value)
            return (
              <div
                key={option.value}
                className={`max-w-56 transition px-4 py-2 bg-gray-50 dark:bg-gray-900 rounded-sm border border-gray-200 dark:border-gray-700 select-none cursor-pointer hover:text-green-800 hover:border-green-300 hover:bg-green-50 ${isSelected ? 'bg-green-50 border-green-300 text-green-800 dark:text-green-200' : ''}`}
                onClick={() => handleOptionClick(option)}
              >
                <div className="flex items-center justify-between gap-1">
                  <span
                    className={`text-sm font-semibold mr-2 text-gray-700 dark:text-gray-200 ${isSelected ? 'text-green-800' : ''}`}
                  >
                    {option.label}
                  </span>
                  <AnimatePresence mode="popLayout">
                    {isSelected ? (
                      <motion.div
                        key="checked"
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        transition={{ duration: 0.1, ease: 'easeInOut' }}
                        className="flex items-center justify-center"
                      >
                        <RemixIcon
                          icon={riCheckboxCircleFill}
                          className="text-green-800 dark:text-green-400"
                          size="sm"
                        />
                      </motion.div>
                    ) : (
                      <motion.div
                        key="unchecked"
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        transition={{ duration: 0.1, ease: 'easeInOut' }}
                        className="flex items-center justify-center"
                      >
                        <RemixIcon
                          icon={riCheckboxBlankCircleLine}
                          className="text-green-800 dark:text-green-400"
                          size="sm"
                        />
                      </motion.div>
                    )}
                  </AnimatePresence>
                </div>
                {option.description && (
                  <p className="text-sm text-gray-600">{option.description}</p>
                )}
              </div>
            )
          })}
        </div>
      </div>
    )
  }

  return (
    <div className={`flex relative my-3 ${className}`} ref={dropdownRef}>
      {/* Add hidden select element for form data */}
      <select
        name={name}
        value={
          multiple ? selectedOptions.map((opt) => opt.value) : (selectedOptions[0]?.value ?? '')
        }
        onChange={onChange}
        multiple={multiple}
        hidden
        aria-hidden="true"
        disabled={disabled}
      >
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>

      <div
        role="combobox"
        aria-expanded={isOpen}
        aria-haspopup="listbox"
        aria-controls="dropdown-list"
        aria-labelledby={`${name}-label`}
        tabIndex={disabled ? -1 : 0}
        className={`${calculatedInputClasses} flex items-center justify-between select-none text-sm`}
        onClick={() => !disabled && setIsOpen(!isOpen)}
        onKeyDown={handleKeyDown}
      >
        <span
          className={`${!selectedOptions.length ? 'text-gray-400 dark:text-gray-500' : 'text-gray-700 dark:text-gray-200'}`}
        >
          {selectedOptions.length
            ? multiple
              ? selectedOptions.map((opt) => opt.label).join(', ')
              : selectedOptions[0].label
            : 'Select an option'}
        </span>
        {optional && selectedOptions.length > 0 && (
          <button
            type="button"
            onClick={handleClearSelection}
            className="absolute transition cursor-pointer text-red-500 hover:text-red-600 dark:hover:text-red-400 top-1/2 -translate-y-1/2 -right-7"
            aria-label="Clear selection"
          >
            <RemixIcon icon={riCloseLine} size="lg" />
          </button>
        )}
        <AnimatePresence initial={false} mode="wait">
          <motion.div
            key={isOpen ? 'up' : 'down'}
            initial={{ opacity: 0, y: isOpen ? 5 : -5 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: isOpen ? -5 : 5 }}
            transition={{ duration: 0.1 }}
          >
            <RemixIcon
              icon={isOpen ? riArrowUpSLine : riArrowDownSLine}
              size="lg"
              className="cursor-pointer"
              aria-hidden="true"
            />
          </motion.div>
        </AnimatePresence>
      </div>

      {isOpen && !disabled && (
        <div
          id="dropdown-list"
          role="listbox"
          aria-multiselectable={multiple}
          className={`absolute left-0 right-0 p-2 flex flex-col gap-1 top-full mt-1 ${backgroundColor} border border-gray-300 dark:border-gray-700 rounded-sm max-h-60 overflow-y-auto z-50`}
        >
          {options.map((option) => {
            const isSelected = selectedOptions.some((opt) => opt.value === option.value)
            return (
              <div
                key={option.value}
                role="option"
                aria-selected={isSelected}
                className={`px-6 py-2 text-sm select-none transition rounded cursor-pointer flex items-center justify-between border border-transparent text-gray-700 hover:bg-green-100 dark:hover:bg-green-800 hover:border-green-200 dark:hover:border-green-700 hover:text-green-700 dark:hover:text-green-300 ${
                  isSelected
                    ? 'bg-green-50 border-green-300 text-green-800 dark:bg-green-700 dark:border-green-400 dark:text-white'
                    : 'dark:text-gray-200'
                }`}
                onClick={() => handleOptionClick(option)}
              >
                <span>{option.label}</span>
                {multiple && (
                  <AnimatePresence mode="popLayout">
                    {isSelected ? (
                      <motion.div
                        key="checked"
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        transition={{ duration: 0.1, ease: 'easeInOut' }}
                        className="flex items-center justify-center"
                      >
                        <RemixIcon icon={riCheckboxCircleFill} className="text-green-800" />
                      </motion.div>
                    ) : (
                      <motion.div
                        key="unchecked"
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        transition={{ duration: 0.1, ease: 'easeInOut' }}
                        className="flex items-center justify-center"
                      >
                        <RemixIcon icon={riCheckboxBlankCircleLine} className="text-green-800" />
                      </motion.div>
                    )}
                  </AnimatePresence>
                )}
              </div>
            )
          })}
        </div>
      )}

      <label id={`${name}-label`} htmlFor={name} className={calculatedLabelClasses}>
        {label}
        {optional && !hideOptionalLabel ? ' (optional)' : ''}
      </label>
    </div>
  )
}

export default FormFieldSelect
