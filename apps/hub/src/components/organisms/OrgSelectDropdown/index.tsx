'use client'

import {
  RemixIcon,
  riAddFill,
  riArrowDropDownFill,
  riArrowDropUpFill,
} from '@mwarnerdotme/react-remixicon'
import { Organization } from '@prisma/client'
import { useClickOutside } from '@react-hooks-library/core'
import { useRouter } from 'next/navigation'
import { FC, useRef, useState } from 'react'
import styles from './index.module.scss'
import CreateOrganizationModal from '../../modals/CreateOrganizationModal'
import Button from '@iotea/libs/frontend/components/atoms/Button'

type Props = {
  organizations: Organization[]
  currentOrg?: Organization
}

const OrgSelectDropdown: FC<Props> = ({ organizations, currentOrg }) => {
  const [isDropDownOpen, setIsDropDownOpen] = useState<boolean>(false)

  const router = useRouter()

  const handleToggleDropDown = () => {
    setIsDropDownOpen((current) => !current)
  }

  const handleChangeOrg = (id: string) => {
    setIsDropDownOpen(false)
    router.push(`/organizations/${id}`)
  }

  const dropdownRef = useRef(null)
  useClickOutside(dropdownRef, () => {
    setIsDropDownOpen(false)
  })

  return (
    <>
      <CreateOrganizationModal />
      <div className={styles.dropdown} ref={dropdownRef}>
        <div className={styles.currentSpace} onClick={handleToggleDropDown}>
          <p className="text-gray-700 dark:text-gray-300 font-semibold">
            {currentOrg?.name ?? 'No organizations found!'}
          </p>
          {isDropDownOpen ? (
            <RemixIcon
              className={`${styles.dropdownToggleIcon} ${styles.open}`}
              icon={riArrowDropUpFill}
              size="xl"
            />
          ) : (
            <RemixIcon className={styles.dropdownToggleIcon} icon={riArrowDropDownFill} size="xl" />
          )}
        </div>
        {isDropDownOpen && (
          <>
            <div className="absolute top-8 -left-4 py-3 rounded-sm border p-2 border-gray-300 dark:border-gray-700 bg-white dark:bg-gray-900 z-10 w-full">
              {organizations.map((org) => {
                const { name: orgName, id: orgId } = org

                if (orgId == currentOrg?.id)
                  return (
                    <div
                      className="my-1 cursor-pointer transition bg-green-100 dark:bg-green-700 text-green-700 dark:text-white border border-green-200 dark:border-green-400 py-1 px-4 rounded-sm"
                      key={orgId}
                      onClick={() => handleChangeOrg(orgId)}
                    >
                      {orgName}
                    </div>
                  )

                return (
                  <div
                    className="my-1 cursor-pointer transition rounded-sm block border border-transparent hover:bg-green-100 dark:hover:bg-green-800 hover:border-green-200 dark:hover:border-green-700 hover:text-green-700 dark:hover:text-green-300 py-1 px-4"
                    key={orgId}
                    onClick={() => handleChangeOrg(orgId)}
                  >
                    {orgName}
                  </div>
                )
              })}
              <hr className="my-2 border-gray-300 dark:border-gray-700" />
              <div className="py-1 px-4">
                <Button
                  variant="underline"
                  modalId="createOrganization"
                  className="text-gray-700 dark:text-gray-300"
                >
                  Create a new org <RemixIcon className="ml-1" icon={riAddFill} />
                </Button>
              </div>
            </div>
          </>
        )}
      </div>
    </>
  )
}

export default OrgSelectDropdown
