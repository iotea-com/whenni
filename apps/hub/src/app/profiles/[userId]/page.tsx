'use client'

import updatePassword from '@iotea/hub/actions/updatePassword'
import useSettingsStore, { Theme } from '@iotea/hub/stores/settingsStore'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import FormFieldText from '@iotea/libs/frontend/components/atoms/FormFieldText'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import { FC, useState } from 'react'
import Container from '@iotea/libs/frontend/components/templates/Container'
import FormFieldSelect from '@iotea/libs/frontend/components/atoms/FormFieldSelect'
import { signout } from '@iotea/hub/actions/auth'

const UserPage: FC = () => {
  const theme = useSettingsStore((state) => state.theme)
  const setTheme = useSettingsStore((state) => state.setTheme)
  const showChannelEditorGrid = useSettingsStore((state) => state.showChannelEditorGrid)
  const setShowChannelEditorGrid = useSettingsStore((state) => state.setShowChannelEditorGrid)
  const showFeedbackButton = useSettingsStore((state) => state.showFeedbackButton)
  const setShowFeedbackButton = useSettingsStore((state) => state.setShowFeedbackButton)

  const [newPassword, setNewPassword] = useState('')

  const handlePasswordFormSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()

    const formData = new FormData(e.currentTarget)
    const newPassword = formData.get('newPassword') as string

    const { error } = await updatePassword({ newPassword })

    if (error) {
      addToast({
        title: 'Failed to update password',
        body: error,
        level: 'error',
      })
    } else {
      addToast({
        title: 'Password updated',
        body: 'Your password has been updated successfully.',
        level: 'success',
      })
    }
  }

  const handleSignout = async () => {
    await signout()
  }

  return (
    <>
      <Container>
        <div className="mt-4 mb-6">
          <h1 className="text-xl">Profile</h1>
          <p>
            <small>Preferences and settings for your account.</small>
          </p>
        </div>
      </Container>
      <Container>
        <section className="mb-4">
          <section
            className="relative bg-gray-50 rounded py-5 px-10 border border-gray-200 dark:bg-gray-900 dark:border-gray-800"
            style={{ boxShadow: '3px 3px 10px 0 rgba(0, 0, 0, 0.03)' }}
          >
            <h2>Theme</h2>
            <div className="flex">
              <FormFieldSelect
                name="theme"
                variant="cards"
                options={[
                  { value: Theme.light, label: 'Light' },
                  { value: Theme.dark, label: 'Dark (beta)' },
                ]}
                label="Choose a theme"
                value={theme}
                onChange={(e) => setTheme(e.target.value as Theme)}
              />
            </div>
            <h2 className="mt-4">Feedback</h2>
            <div className="flex">
              <FormFieldSelect
                name="showFeedbackCheckbox"
                variant="cards"
                options={[{ value: 'true', label: 'Show' }]}
                optional
                hideOptionalLabel
                label="Show feedback button"
                value={showFeedbackButton.toString()}
                onChange={(e) => setShowFeedbackButton(e.target.value === 'true')}
              />
            </div>
          </section>
        </section>
      </Container>
      <Container className="mb-4">
        <section
          className="relative bg-gray-50 rounded py-5 px-10 border border-gray-200 dark:bg-gray-900 dark:border-gray-800"
          style={{ boxShadow: '3px 3px 10px 0 rgba(0, 0, 0, 0.03)' }}
        >
          <h2>Channel Editor</h2>
          <div className="flex">
            <FormFieldSelect
              name="showGridCheckbox"
              variant="cards"
              options={[{ value: 'true', label: 'Show' }]}
              optional
              hideOptionalLabel
              label="Show grid"
              value={showChannelEditorGrid.toString()}
              onChange={(e) => setShowChannelEditorGrid(e.target.value === 'true')}
            />
          </div>
        </section>
      </Container>
      <Container className="mb-4">
        <section
          className="relative bg-gray-50 rounded py-5 px-10 border border-gray-200 dark:bg-gray-900 dark:border-gray-800"
          style={{ boxShadow: '3px 3px 10px 0 rgba(0, 0, 0, 0.03)' }}
        >
          <h2>Update Password</h2>
          <form onSubmit={handlePasswordFormSubmit} className="w-96">
            <FormFieldText
              name="newPassword"
              label="New Password"
              inputType="password"
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
            />
            <Button
              className="transition"
              type="submit"
              text="Change password"
              disabled={!newPassword}
            />
          </form>
          <h2 className="mt-4">Sign out</h2>
          <Button text="Sign out" onClick={() => handleSignout()} />
        </section>
      </Container>
    </>
  )
}

export default UserPage
