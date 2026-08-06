import Container from '../../templates/Container'
import Image from 'next/image'
import Link from 'next/link'

import {
  RemixIcon,
  riBlueskyFill,
  riCupFill,
  riDiscordFill,
  riGithubFill,
  riHeartFill,
  riLinkedinFill,
  riTwitterXFill,
  riYoutubeFill,
} from '@mwarnerdotme/react-remixicon'
import MailchimpNewsletterForm from './MailchimpNewsletterForm'

const Footer = async () => {
  return (
    <footer className="pt-24 pb-32 bg-gray-100 z-10">
      <Container>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 md:gap-0">
          <section className="flex flex-col">
            <Image
              src="/img/logos/mark-cutout-secondary.png"
              alt="GRUENT logo"
              width={30}
              height={30}
            />
            <div className="grow hidden md:block" />
            <ul className="flex gap-4 mt-2 md:mt-0">
              <li>
                <Link href="https://discord.gg/prAJjW426d" target="_blank">
                  <RemixIcon icon={riDiscordFill} />
                </Link>
              </li>
              <li>
                <Link href="https://x.com/gruent_com" target="_blank">
                  <RemixIcon icon={riTwitterXFill} />
                </Link>
              </li>
              <li>
                <Link href="https://bsky.app/profile/gruent.bsky.social" target="_blank">
                  <RemixIcon icon={riBlueskyFill} />
                </Link>
              </li>
              <li>
                <Link href="https://github.com/ongruent" target="_blank">
                  <RemixIcon icon={riGithubFill} />
                </Link>
              </li>
              <li>
                <Link href="https://linkedin.com/company/gruent" target="_blank">
                  <RemixIcon icon={riLinkedinFill} />
                </Link>
              </li>
              <li>
                <Link href="https://youtube.com/@gruent" target="_blank">
                  <RemixIcon icon={riYoutubeFill} />
                </Link>
              </li>
            </ul>
            <p className="mt-2 text-gray-500">&copy; {new Date().getFullYear()} GRUENT, Inc.</p>
            <p className="text-gray-400 text-sm">
              Made with <RemixIcon className="text-red-700" icon={riHeartFill} /> and plenty of{' '}
              <RemixIcon className="text-green-700" icon={riCupFill} />
            </p>
          </section>
          <section>
            <h3 className="mt-4 mb-1">Resources</h3>
            <ul>
              <li>
                <Link href="/blog">Blog</Link>
              </li>
              <li>
                <Link href="/docs">Docs</Link>
              </li>
              <li>
                <Link href="/pricing">Pricing</Link>
              </li>
              <li>
                <Link href="mailto:careers@gruent.com?subject=I'd like to work at GRUENT">
                  Careers
                </Link>
              </li>
              {/* <li>
                <Link href="/blog">Roadmap</Link>
              </li> */}
            </ul>
          </section>
          <section>
            <h3 className="mt-4 mb-1">Newsletter</h3>
            <MailchimpNewsletterForm />
          </section>
          {/* <section>
            <h3>Products</h3>
            <ul>
              <li>
                <Link href="/products/grid">Grid</Link>
              </li>
              <li>
                <Link href="/products/grid/on-prem">On-prem</Link>
              </li>
              <li>
                <Link href="/products/data-logger">Data Logger</Link>
              </li>
              <li>
                <Link href="/products/carrier-board">Carrier Board</Link>
              </li>
            </ul>
          </section>
          <section>
            <h3>Developers</h3>
            <ul>
              <li>
                <Link href="/developers/documentation">Documentation</Link>
              </li>
              <li>
                <Link href="/status">API Status</Link>
              </li>
            </ul>
            <h3>Company</h3>
            <ul>
              <li>
                <Link href="/careers">Careers</Link>
              </li>
              <li>
                <Link href="/about">About</Link>
              </li>
              <li>
                <Link href="/about#team">Team</Link>
              </li>
              <li>
                <Link href="/partners">Become a Partner</Link>
              </li>
            </ul>
          </section>
          <section>
            <h3>Resources</h3>
            <ul>
              <li>
                <Link href="/support">Support</Link>
              </li>
              <li>
                <Link href="/blog">Blog</Link>
              </li>
              <li>
                <Link href="/pricing">Pricing</Link>
              </li>
              <li>
                <Link href="/blog">Roadmap</Link>
              </li>
              <li>
                <Link href="/use-cases">Use Cases</Link>
              </li>
              <li>
                <Link href="/legal">Privacy & Terms</Link>
              </li>
              <li>
                <Link href="/sitemap">Sitemap</Link>
              </li>
            </ul>
          </section> */}
        </div>
      </Container>
    </footer>
  )
}

export default Footer
