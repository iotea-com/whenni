'use client'

import { handleUpdatePolicy } from '@gruent/hub/actions/policies'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import FormFieldText from '@gruent/libs/frontend/components/atoms/FormFieldText'
import { addToast } from '@gruent/libs/frontend/hooks/useToast'
import { Policy } from '@gruent/libs/gruent-js/src/policies'
import { RemixIcon, riCloseLine } from '@mwarnerdotme/react-remixicon'
import { FC, useCallback, useMemo, useState } from 'react'
import gruentClient from '@gruent/hub/lib/gruent'
import FormFieldSelect from '@gruent/libs/frontend/components/atoms/FormFieldSelect'
import useAuth from '@gruent/hub/hooks/useAuth'

type Props = {
  policy: Policy
  revoke: boolean
  spaceId: string
  certificateId: string
  className?: string
}

const PolicySection: FC<Props> = ({ policy, revoke, spaceId, certificateId, className }) => {
  const { accessToken } = useAuth()

  const [isRevokeLoading, setIsRevokeLoading] = useState(false)
  const [allowedSubscriptionTopics, setAllowedSubscriptionTopics] = useState(
    policy.allowedSubscriptionTopics,
  )
  const [allowedPublishTopics, setAllowedPublishTopics] = useState(policy.allowedPublishTopics)

  const showSaveChangesButton = useMemo(() => {
    if (
      JSON.stringify(policy.allowedSubscriptionTopics) !== JSON.stringify(allowedSubscriptionTopics)
    )
      return true

    if (JSON.stringify(policy.allowedPublishTopics) !== JSON.stringify(allowedPublishTopics))
      return true

    return false
  }, [allowedPublishTopics, allowedSubscriptionTopics, policy])

  const handleRevokeCheckbox = useCallback(
    async (revoke: boolean) => {
      setIsRevokeLoading(true)

      const { error } = await handleUpdatePolicy(spaceId, certificateId, policy, revoke)

      if (error) {
        addToast({
          title: 'Could not revoke the certificate',
          body: error,
          level: 'error',
        })
        setIsRevokeLoading(false)
        return
      }

      addToast({
        title: revoke
          ? 'Successfully revoked the certificate'
          : 'Successfully enabled the certificate',
        body: revoke
          ? 'This MQTT client will no longer be able to subscribe to any topics or publish any messages.'
          : 'This MQTT client can now subscribe to topics and publish messages on allowed topics.',
        level: 'success',
      })
      setIsRevokeLoading(false)
    },
    [certificateId, policy, spaceId],
  )

  const handleSavePolicy = useCallback(async () => {
    const newPolicy: Policy = {
      allowedSubscriptionTopics,
      allowedPublishTopics,
    }

    const { error } = await handleUpdatePolicy(spaceId, certificateId, newPolicy, revoke)

    if (error) {
      addToast({
        title: 'Could not update the policy',
        body: error,
        level: 'error',
      })
      setIsRevokeLoading(false)
      return
    }

    addToast({
      title: 'Successfully updated the policy',
      body: 'The new policy has been saved and will be implemented in the next minute.',
      level: 'success',
    })
  }, [certificateId, spaceId, allowedPublishTopics, allowedSubscriptionTopics, revoke])

  const handleDownloadBundleClick = useCallback(async () => {
    // Get current access token
    if (!accessToken) return

    // Retrieve certificate bundle
    const { data: certificateResponseData, errors } = await gruentClient(
      accessToken,
    ).certificates.get(spaceId, certificateId)

    if (errors && errors.length > 0) {
      addToast({
        title: 'Could not retrieve the certificate bundle',
        body: errors[0],
        level: 'error',
      })
      setIsRevokeLoading(false)
      return
    }

    if (!certificateResponseData) {
      addToast({
        title: 'Could not retrieve the certificate bundle',
        body: 'Invalid response data.',
        level: 'error',
      })
      setIsRevokeLoading(false)
      return
    }

    // Show download window
    const { blob, filename } = certificateResponseData
    const certificateBundleZip = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = certificateBundleZip
    link.setAttribute('download', filename)
    document.body.appendChild(link)
    link.click()

    // Clean up
    link.remove()
    window.URL.revokeObjectURL(certificateBundleZip)
  }, [accessToken, certificateId, spaceId])

  const handleEditTopic = useCallback(
    (action: 'subscribe' | 'publish', index: number, topic: string) => {
      if (!topic.startsWith(`${spaceId}/`)) {
        addToast({
          title: 'Cannot change topic',
          body: `The topic could not be changed to "${topic}" because the topic must start with the space ID. For example, if the space ID is "s123", the topic must start with "s123/".`,
          level: 'error',
        })
        return
      }

      if (action === 'subscribe') {
        setAllowedSubscriptionTopics((current) => {
          const updated = [...current]
          updated[index] = topic
          return updated
        })
      }

      if (action === 'publish') {
        setAllowedPublishTopics((current) => {
          const updated = [...current]
          updated[index] = topic
          return updated
        })
      }
    },
    [spaceId],
  )

  const handleRemoveTopic = useCallback((action: 'subscribe' | 'publish', index: number) => {
    if (action === 'subscribe') {
      setAllowedSubscriptionTopics((current) => {
        const updated = [...current]
        updated.splice(index, 1)
        return updated
      })
    }

    if (action === 'publish') {
      setAllowedPublishTopics((current) => {
        const updated = [...current]
        updated.splice(index, 1)
        return updated
      })
    }
  }, [])

  return (
    <section className={className}>
      <h2>Certificate</h2>
      <Button text="Download bundle" onClick={handleDownloadBundleClick} />
      <h2 className="mt-4">Policy</h2>
      <h3 className="mt-2 text-sm">Allowed Subscription Topics</h3>
      {allowedSubscriptionTopics.map((topic, i) => {
        return (
          <div key={`topic${i + 1}`} className="relative">
            <FormFieldText
              name={`topic${i + 1}`}
              label={`Topic ${i + 1}`}
              value={topic}
              onChange={(e) => handleEditTopic('subscribe', i, e.target.value)}
            />
            <RemixIcon
              icon={riCloseLine}
              className="absolute -right-5 top-1/2 -translate-y-1/2 cursor-pointer"
              onClick={() => handleRemoveTopic('subscribe', i)}
            />
          </div>
        )
      })}
      <Button
        text="Add subscription topic"
        onClick={() => setAllowedSubscriptionTopics((current) => [...current, `${spaceId}/`])}
      />
      <h3 className="mt-2 text-sm">Allowed Publish Topics</h3>
      {allowedPublishTopics.map((topic, i) => {
        return (
          <div key={`topic${i + 1}`} className="relative">
            <FormFieldText
              name={`topic${i + 1}`}
              label={`Topic ${i + 1}`}
              value={topic}
              onChange={(e) => handleEditTopic('publish', i, e.target.value)}
            />
            <RemixIcon
              icon={riCloseLine}
              className="absolute -right-5 top-1/2 -translate-y-1/2 cursor-pointer"
              onClick={() => handleRemoveTopic('publish', i)}
            />
          </div>
        )
      })}
      <Button
        text="Add publish topic"
        onClick={() => setAllowedPublishTopics((current) => [...current, `${spaceId}/`])}
      />
      {showSaveChangesButton && (
        <div>
          <Button className="mt-4" text="Save changes" onClick={handleSavePolicy} />
        </div>
      )}
      <h2 className="mt-4">Revoke certificate</h2>
      <FormFieldSelect
        name="revoke"
        variant="cards"
        options={[{ value: 'true', label: 'Revoke' }]}
        optional
        hideLabel
        disabled={isRevokeLoading}
        label="Revoke"
        value={revoke.toString()}
        onChange={(e) => handleRevokeCheckbox(e.target.value === 'true')}
      />
    </section>
  )
}

export default PolicySection
