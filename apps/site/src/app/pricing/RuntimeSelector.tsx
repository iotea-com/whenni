'use client'

import FormFieldSelect from '@gruent/libs/frontend/components/atoms/FormFieldSelect'

type RuntimeSelectorProps = {
  runtime: 'sm' | 'md' | 'lg'
}

const runtimeSelector = ({ runtime }: RuntimeSelectorProps) => {
  return (
    <FormFieldSelect
      name="runtime"
      options={[
        { label: 'Small', value: 'sm' },
        { label: 'Medium', value: 'md' },
        { label: 'Large', value: 'lg' },
      ]}
      label="Runtime size"
      value={runtime}
      className="max-w-48 mx-auto"
      onChange={(e) => {
        window.location.search = `?runtime=${e.target.value}`
      }}
    />
  )
}

export default runtimeSelector
