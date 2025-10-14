'use server'

import { btoa } from 'buffer'

const addNewsletterMember = async (email: string) => {
  const apiKey = btoa(`any:${process.env.MAILCHIMP_API_KEY}`)

  const addMemberResult = await fetch(
    `https://api.mailchimp.com/3.0/lists/${process.env.MAILCHIMP_LIST_ID}/members`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Basic ${apiKey}`,
      },
      body: JSON.stringify({ email_address: email, status: 'pending' }),
    },
  )

  switch (addMemberResult.status) {
    case 200:
      return {
        message: `Thanks! A confirmation email is on it's way to your inbox.`,
      }
    case 400:
      return { message: `Good news! You're already subscribed.` }
    default:
      return { message: 'Could not use this email. Please try again.' }
  }
}

export default addNewsletterMember
