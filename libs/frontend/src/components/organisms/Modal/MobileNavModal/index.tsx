'use client'

import Modal from '@iotea/libs/frontend/components/organisms/Modal'
import { closeModal } from '@iotea/libs/frontend/hooks/useModal'
import { RemixIcon, riCloseFill } from '@mwarnerdotme/react-remixicon'
import Image from 'next/image'
import Link from 'next/link'

const MobileNavModal = () => {
  return (
    <Modal id="mobileNav" showAccept={false} className="relative">
      <RemixIcon
        icon={riCloseFill}
        size="lg"
        className="absolute top-8 right-5"
        onClick={closeModal}
      />
      <Image
        src="/img/logos/mark-cutout-secondary.png"
        alt="IOTEA logo"
        width={30}
        height={30}
        className="mb-4"
      />
      <ul className="flex flex-col gap-2">
        <li>
          <Link href="/" onClick={closeModal}>
            Home
          </Link>
        </li>
        <li>
          <Link href="/blog" onClick={closeModal}>
            Blog
          </Link>
        </li>
        <li>
          <Link href="/docs" onClick={closeModal}>
            Docs
          </Link>
        </li>
        <li>
          <Link href="/pricing" onClick={closeModal}>
            Pricing
          </Link>
        </li>
      </ul>
    </Modal>
  )
}

export default MobileNavModal
