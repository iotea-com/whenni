'use client'

import {
  RemixIcon,
  riAddFill,
  riArrowDropDownLine,
  riArrowDropUpLine,
} from '@mwarnerdotme/react-remixicon'
import { Organization, Space } from '@prisma/client'
import { useClickOutside } from '@react-hooks-library/core'
import { FC, useMemo, useRef, useState } from 'react'
import styles from './index.module.scss'
import CreateSpaceModal from '@iotea/hub/components/modals/CreateSpaceModal'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import Link from 'next/link'

type Props = {
  organization?: Organization & { spaces: Space[] }
  currentSpaceId: string
}

const SpaceSelectDropDown: FC<Props> = ({ organization, currentSpaceId }) => {
  const [isDropDownOpen, setIsDropDownOpen] = useState<boolean>(false)

  const handleToggleDropDown = () => {
    setIsDropDownOpen((current) => !current)
  }

  const dropdownRef = useRef(null)
  useClickOutside(dropdownRef, () => {
    setIsDropDownOpen(false)
  })

  const currentSpace = useMemo(() => {
    if (!organization?.spaces || organization.spaces.length <= 0) return

    // get the current space from the users array
    const cs = organization.spaces.filter((space) => {
      if (space.id === currentSpaceId) return true
      return false
    })

    if (cs.length > 0) return cs[0]
    return organization.spaces[0]
  }, [organization?.spaces, currentSpaceId])

  if (!currentSpace || !organization?.spaces || organization.spaces.length <= 0)
    return <p>No spaces were found!</p>

  if (!organization) return <p className="text-gray-700 dark:text-gray-300 font-semibold">Error</p>

  return (
    <>
      <CreateSpaceModal orgId={organization.id} />
      <div className={styles.dropdown} ref={dropdownRef}>
        <div className={styles.currentSpace} onClick={handleToggleDropDown}>
          <p className="text-gray-700 dark:text-gray-300 font-semibold">{currentSpace.name}</p>
          {isDropDownOpen ? (
            <RemixIcon
              className={`${styles.dropdownToggleIcon} ${styles.open}`}
              icon={riArrowDropUpLine}
              size="xl"
            />
          ) : (
            <RemixIcon className={styles.dropdownToggleIcon} icon={riArrowDropDownLine} size="xl" />
          )}
        </div>
        {isDropDownOpen && (
          <>
            <div className="absolute top-8 -left-4 py-3 rounded-sm border p-2 border-gray-300 dark:border-gray-700 bg-white dark:bg-gray-900 z-10 w-full">
              {organization.spaces.map((space) => {
                const { name: spaceName, id: spaceId } = space

                if (spaceId == currentSpaceId)
                  return (
                    <div
                      className="my-1 cursor-pointer transition bg-green-100 dark:bg-green-700 text-green-700 dark:text-white border border-green-200 dark:border-green-400 py-1 px-4 rounded-sm"
                      key={spaceId}
                    >
                      {spaceName}
                    </div>
                  )

                return (
                  <Link
                    key={spaceId}
                    className="my-1 cursor-pointer transition rounded-sm block border border-transparent hover:bg-green-100 dark:hover:bg-green-800 hover:border-green-200 dark:hover:border-green-700 hover:text-green-700 dark:hover:text-green-300 py-1 px-4"
                    href={`/organizations/${organization.id}/spaces/${spaceId}`}
                  >
                    {spaceName}
                  </Link>
                )
              })}
              <hr className="my-2 border-gray-300 dark:border-gray-700" />
              <div className="py-1 px-4">
                <Button
                  variant="underline"
                  modalId="createSpace"
                  className="text-gray-700 dark:text-gray-300"
                >
                  Create a new space <RemixIcon className="ml-1" icon={riAddFill} />
                </Button>
              </div>
            </div>
          </>
        )}
      </div>
    </>
  )
}

export default SpaceSelectDropDown
