'use client'

import { FC } from 'react'
import Container from '../../templates/Container'
import Image from 'next/image'
import Link from 'next/link'
import Button from '../../atoms/Button'
import { RemixIcon, riMenu2Fill } from '@mwarnerdotme/react-remixicon'
import { openModal } from '@gruent/libs/frontend/hooks/useModal'

const Navbar: FC = () => {
  return (
    <div className="py-4 absolute top-0 left-0 w-full z-30 bg-white border-b border-gray-100">
      <Container>
        <div className="flex items-center">
          <Link href="/">
            <Image
              src="/img/logos/wordmark-secondary.png"
              alt="GRUENT logo"
              width={100}
              height={100}
            />
          </Link>
          <nav className="hidden md:flex justify-between items-center w-full">
            <ul className="ml-8 flex gap-8 text-sm">
              <li>
                <Link href="/">Home</Link>
              </li>
              <li>
                <Link href="/blog">Blog</Link>
              </li>
              <li>
                <Link href="/docs">Docs</Link>
              </li>
              <li>
                <Link href="/pricing">Pricing</Link>
              </li>
              {/* <li>
                <Link href="/products">Products</Link>
              </li>
              <li>
                <Link href="/developers">Developers</Link>
              </li>
              <li>
                <Link href="/resources">Resources</Link>
              </li>
              <li>
                <Link href="/roadmap">Roadmap</Link>
              </li> */}
            </ul>
            <ul className="flex gap-4">
              {/* <li>
                <Button
                  text="Sign In"
                  href="https://app.gruent.com/sign-in"
                  variant="underline"
                />
              </li>
              */}
              <li className="font-bold">
                <Button text="Get Started" href="#" onClick={() => openModal('betaSignup')} />
              </li>
            </ul>
          </nav>
          <div className="grow md:hidden" />
          <div
            className="block md:hidden w-9 h-9 flex justify-center items-center bg-secondary cursor-pointer transition rounded-full text-white"
            onClick={() => openModal('mobileNav')}
          >
            <RemixIcon icon={riMenu2Fill} />
          </div>
        </div>
      </Container>
    </div>
  )
}

export default Navbar
