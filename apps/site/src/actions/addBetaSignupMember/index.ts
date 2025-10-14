'use server'

import { btoa } from 'buffer'

const addBetaSignupMember = async (email: string) => {
  const apiKey = btoa(`any:${process.env.MAILCHIMP_API_KEY}`)
  const baseUrl = `https://api.mailchimp.com/3.0/lists/${process.env.MAILCHIMP_LIST_ID}`

  // Add member to list
  const addMemberResult = await fetch(`${baseUrl}/members`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Basic ${apiKey}`,
    },
    body: JSON.stringify({
      email_address: email,
      status: 'subscribed',
      tags: ['Open Beta User'],
    }),
  })

  if (addMemberResult.status !== 200 && addMemberResult.status !== 400) {
    return { message: 'Could not use this email. Please try again.' }
  }

  // Return appropriate message
  switch (addMemberResult.status) {
    case 200:
      return {
        message: `Thanks! A confirmation email is on its way to your inbox.`,
      }
    case 400:
      return { message: `Good news! You're already on the list.` }
    default:
      return { message: 'Could not use this email. Please try again.' }
  }
}

export default addBetaSignupMember
