'use client'

import { ChangeEventHandler, FC, useMemo, useState } from 'react'
import FormFieldSelect from '../../atoms/FormFieldSelect'

const getStepSize = (value: number, maximumRequestsPerSecond: number) => {
  // Calculate reasonable step sizes based on max requests
  const maxStep = Math.max(1, Math.floor(maximumRequestsPerSecond / 100)) // Ensure at least 100 steps

  if (value < 1000) return Math.min(100, maxStep) // Steps of 100 or less up to 1K
  if (value < 10000) return Math.min(1000, maxStep) // Steps of 1K or less up to 10K
  if (value < 100000) return Math.min(10000, maxStep) // Steps of 10K or less up to 100K
  if (value < 1000000) return Math.min(100000, maxStep) // Steps of 100K or less up to 1M
  if (value < 10000000) return Math.min(1000000, maxStep) // Steps of 1M or less up to 10M
  return Math.min(10000000, maxStep) // Steps of 10M or less up to 1B
}

type Props = {
  runtime: 'sm' | 'md' | 'lg'
  baseCost: number
}

const ChannelsCostCalculator: FC<Props> = ({
  runtime: initialRuntime,
  baseCost: initialBaseCost,
}: Props) => {
  const [numChannelExecutions, setNumChannelExecutions] = useState<number>(0)
  const [ingressDataKb, setIngressDataKb] = useState<number>(0)
  const [egressDataKb, setEgressDataKb] = useState<number>(0)
  const [runtime, setRuntime] = useState<'sm' | 'md' | 'lg'>(initialRuntime)
  const [baseCost, setBaseCost] = useState<number>(initialBaseCost)

  const handleNumChannelExecutionsChange: ChangeEventHandler<HTMLInputElement> = (event) => {
    setNumChannelExecutions(Number(event.target.value))
  }

  const maximumRequestsPerSecond = useMemo(() => {
    const oneRequestPerSecond = 2592000 // 30 days * 24 hours * 60 minutes * 60 seconds

    const maximumRequestsPerSecond = (() => {
      switch (runtime) {
        case 'sm':
          return 50
        case 'md':
          return 200
        case 'lg':
          return 500
      }
    })()

    return oneRequestPerSecond * maximumRequestsPerSecond
  }, [runtime])

  const cost = useMemo(() => {
    return calculateCost(
      runtime,
      baseCost,
      numChannelExecutions,
      ingressDataKb * 1e3,
      egressDataKb * 1e3,
    )
  }, [runtime, baseCost, numChannelExecutions, ingressDataKb, egressDataKb])

  return (
    <div>
      <div className="mt-2">
        <div className="flex flex-col">
          <FormFieldSelect
            name="runtime"
            options={[
              { label: 'Small', value: 'sm' },
              { label: 'Medium', value: 'md' },
              { label: 'Large', value: 'lg' },
            ]}
            label="Runtime size"
            value={runtime}
            onChange={(e) => {
              const newRuntime = e.target.value as 'sm' | 'md' | 'lg'
              const newBaseCost = (() => {
                switch (newRuntime) {
                  case 'sm':
                    return 5
                  case 'md':
                    return 10
                  case 'lg':
                    return 15
                }
              })()
              setRuntime(newRuntime)
              setBaseCost(newBaseCost)
            }}
          />
          <p className="text-xs text-gray-500 mb-4">
            Estimated maximum number of requests/executions per second for this runtime:{' '}
            {maximumRequestsPerSecond.toLocaleString()}
          </p>
          <label htmlFor="numChannelExecutions">
            How many times will you run the channel every month?
          </label>
          <input
            type="range"
            id="numChannelExecutions"
            min={0}
            max={maximumRequestsPerSecond}
            step={getStepSize(numChannelExecutions, maximumRequestsPerSecond)}
            value={numChannelExecutions}
            onChange={handleNumChannelExecutionsChange}
          />
          <span>{numChannelExecutions.toLocaleString()}</span>
          <label htmlFor="ingressDataBytes">
            How much data (in kilobytes) will you transfer in per channel execution?
          </label>
          <input
            type="range"
            id="ingressDataBytes"
            min={1}
            max={1024}
            value={ingressDataKb}
            onChange={(e) => setIngressDataKb(Number(e.target.value))}
          />
          <span>{ingressDataKb.toLocaleString()}kb</span>
          <label htmlFor="egressDataBytes">
            How much data (in kilobytes) will you transfer out per channel execution?
          </label>
          <input
            type="range"
            id="egressDataBytes"
            min={1}
            max={1024}
            value={egressDataKb}
            onChange={(e) => setEgressDataKb(Number(e.target.value))}
          />
          <span>{egressDataKb.toLocaleString()}kb</span>
        </div>
      </div>

      <div className="mt-2">
        <p className="font-semibold">Estimated cost: ${cost}/mo + sales tax</p>
        <p className="text-xs text-gray-500">Pricing subject to change throughout beta.</p>
      </div>
    </div>
  )
}

const runtimeTiers = {
  sm: [
    { limit: 1000, cost: 0.5 },
    { limit: 10000, cost: 0.25 },
    { limit: 100000, cost: 0.05 },
    { limit: 1000000, cost: 0.01 },
    { limit: 10000000, cost: 0.005 },
    { limit: 1000000000, cost: 0.001 },
  ],
  md: [
    { limit: 1000, cost: 0.25 },
    { limit: 10000, cost: 0.1 },
    { limit: 100000, cost: 0.025 },
    { limit: 1000000, cost: 0.01 },
    { limit: 10000000, cost: 0.005 },
    { limit: 1000000000, cost: 0.001 },
  ],
  lg: [
    { limit: 1000, cost: 0.1 },
    { limit: 10000, cost: 0.01 },
    { limit: 100000, cost: 0.01 },
    { limit: 1000000, cost: 0.01 },
    { limit: 10000000, cost: 0.005 },
    { limit: 1000000000, cost: 0.001 },
  ],
}

const calculateCost = (
  runtime: 'sm' | 'md' | 'lg',
  baseCost: number,
  numChannelExecutions: number,
  ingressDataBytes: number,
  egressDataBytes: number,
) => {
  const ingressCost = ((ingressDataBytes * numChannelExecutions) / 1e9) * 0.1
  const egressCost = ((egressDataBytes * numChannelExecutions) / 1e9) * 0.15

  let remainingExecutions = numChannelExecutions
  let totalCost = ingressCost + egressCost

  // Get the last tier's cost for any executions beyond the defined tiers
  const lastTier = runtimeTiers[runtime][runtimeTiers[runtime].length - 1]

  for (let i = 0; i < runtimeTiers[runtime].length; i++) {
    const currentTier = runtimeTiers[runtime][i]
    const previousLimit = i === 0 ? 0 : runtimeTiers[runtime][i - 1].limit
    const executionsInTier = Math.min(remainingExecutions, currentTier.limit - previousLimit)

    if (executionsInTier <= 0) break

    const tierCost = (executionsInTier / 1000) * currentTier.cost
    totalCost += tierCost
    remainingExecutions -= executionsInTier
  }

  // Handle any remaining executions using the last tier's rate
  if (remainingExecutions > 0) {
    totalCost += (remainingExecutions / 1000) * lastTier.cost
  }

  if (totalCost < baseCost) return baseCost.toFixed(2)

  return totalCost.toFixed(2)
}

export default ChannelsCostCalculator
