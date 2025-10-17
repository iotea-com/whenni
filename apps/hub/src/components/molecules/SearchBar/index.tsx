'use client'

import { RemixIcon, riCloseLine, riSearchLine } from '@mwarnerdotme/react-remixicon'
import { useRouter } from 'next/navigation'
import { FC, useEffect, useState } from 'react'

type Props = {
  paramKey?: string
  className?: string
  initialFilter?: string
  placeholder?: string
  onChange?: (filter: string) => void
}

const SearchBar: FC<Props> = ({
  paramKey = 'q',
  className,
  initialFilter,
  placeholder,
  onChange,
}) => {
  const router = useRouter()

  const [search, setSearch] = useState<string>(initialFilter ?? '')

  useEffect(() => {
    if (search) {
      const newLocation = new URL(window.location.toString())
      newLocation.searchParams.set(paramKey, search)
      router.replace(newLocation.toString())
    } else {
      const newLocation = new URL(window.location.toString())
      newLocation.searchParams.delete(paramKey)
      router.replace(newLocation.toString())
    }

    if (onChange) onChange(search)
  }, [search, router, paramKey, onChange])

  return (
    <div className="relative">
      <RemixIcon
        icon={riSearchLine}
        className="absolute left-1 top-1/2 -translate-y-1/2 text-gray-500"
      />
      {search && (
        <RemixIcon
          icon={riCloseLine}
          className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-500 cursor-pointer"
          onClick={() => setSearch('')}
        />
      )}
      <input
        type="text"
        onChange={(e) => setSearch(e.target.value)}
        value={search}
        placeholder={placeholder}
        className={`${className} w-full px-8 py-1 border-b bg-transparent transition outline-hidden border-gray-200 hover:border-gray-300 focus:border-gray-700 text-sm text-gray-700`}
      />
    </div>
  )
}

export default SearchBar
