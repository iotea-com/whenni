import useChannelEditorStore from '@gruent/hub/stores/channelEditorStore'
import FormFieldText from '@gruent/libs/frontend/components/atoms/FormFieldText'
import { ChannelNode } from '@gruent/libs/engine/nodes/v1'
import { TimerSourceConfig } from '@gruent/libs/engine/nodes/v1/src/source/timer'
import { FC, useEffect, useState, useMemo } from 'react'
import { RemixIcon, riArrowDownSFill } from '@mwarnerdotme/react-remixicon'

import FormFieldSelect from '@gruent/libs/frontend/components/atoms/FormFieldSelect'

const TimerSourceOptions: FC = () => {
  const upsertNode = useChannelEditorStore((state) => state.upsertNode)
  const currentNode = useChannelEditorStore(
    (state) => state.currentNode as ChannelNode<TimerSourceConfig>,
  )

  const [showAdvanced, setShowAdvanced] = useState(false)

  const [seconds, setSeconds] = useState<string>('*')
  const [minutes, setMinutes] = useState<string>('*')
  const [hours, setHours] = useState<string>('*')
  const [dayOfMonth, setDayOfMonth] = useState<string>('*')
  const [month, _setMonth] = useState<string>('*')
  const [dayOfWeek, setDayOfWeek] = useState<string>('*')
  const [secondsExact, setSecondsExact] = useState<boolean>(false)
  const [minutesExact, setMinutesExact] = useState<boolean>(false)
  const [hoursExact, setHoursExact] = useState<boolean>(false)
  const [dayOfMonthExact, setDayOfMonthExact] = useState<boolean>(false)
  const [monthExact, _setMonthExact] = useState<boolean>(false)
  const [dayOfWeekExact, setDayOfWeekExact] = useState<boolean>(false)

  const [now, setNow] = useState<Date>(new Date())

  const nextExecutionTimes = useMemo(() => {
    // Start with current time
    let currentDate = new Date(now.getTime())
    const executionTimes: Date[] = []

    // Parse all values to handle both exact and interval modes
    const parseField = (value: string, isExact: boolean) => {
      if (value === '*') return null
      return isExact ? parseInt(value) : parseInt(value)
    }

    const secondVal = parseField(seconds, secondsExact)
    const minuteVal = parseField(minutes, minutesExact)
    const hourVal = parseField(hours, hoursExact)
    const dayOfMonthVal = parseField(dayOfMonth, dayOfMonthExact)
    const monthVal = parseField(month, monthExact)
    const dayOfWeekVal = parseField(dayOfWeek, dayOfWeekExact)

    // Function to check if a date matches our cron criteria
    const matchesCriteria = (date: Date) => {
      // For exact values, check if they match
      if (secondVal !== null && secondsExact && date.getSeconds() !== secondVal) return false
      if (minuteVal !== null && minutesExact && date.getMinutes() !== minuteVal) return false
      if (hourVal !== null && hoursExact && date.getHours() !== hourVal) return false
      if (dayOfMonthVal !== null && dayOfMonthExact && date.getDate() !== dayOfMonthVal)
        return false
      if (monthVal !== null && monthExact && date.getMonth() !== monthVal - 1) return false
      if (dayOfWeekVal !== null && dayOfWeekExact && date.getDay() !== dayOfWeekVal % 7)
        return false

      // For interval values, check if they're divisible
      if (secondVal !== null && !secondsExact && date.getSeconds() % secondVal !== 0) return false
      if (minuteVal !== null && !minutesExact && date.getMinutes() % minuteVal !== 0) return false
      if (hourVal !== null && !hoursExact && date.getHours() % hourVal !== 0) return false
      // Day of month and month intervals are more complex and not fully implemented here

      return true
    }

    // Find the next 3 execution times
    const maxIterations = 365 * 24 * 60 * 60 // 1 year max
    let iterations = 0

    // Add 1 second to start from the next second
    currentDate.setSeconds(currentDate.getSeconds() + 1)

    while (executionTimes.length < 3 && iterations < maxIterations) {
      if (matchesCriteria(currentDate)) {
        // Clone the date to avoid reference issues
        executionTimes.push(new Date(currentDate.getTime()))
      }
      currentDate.setSeconds(currentDate.getSeconds() + 1)
      iterations++
    }

    // If we couldn't find enough matches, return what we have
    if (iterations >= maxIterations && executionTimes.length === 0) {
      return ['Could not determine next execution times']
    }

    // Format the dates
    return executionTimes.map((date) => date.toLocaleString())
  }, [
    now,
    seconds,
    minutes,
    hours,
    dayOfMonth,
    month,
    dayOfWeek,
    secondsExact,
    minutesExact,
    hoursExact,
    dayOfMonthExact,
    monthExact,
    dayOfWeekExact,
  ])

  // Load in default values
  useEffect(() => {
    const currentConfig = { ...currentNode.metadata.config }

    const parts = currentConfig.cronExpression.split(' ')

    const parseCrontabField = (value: string, min: number, max: number) => {
      if (value === '*') return { isExact: false, num: 1 }
      let num = parseInt(value.match(/\d+/)?.[0] ?? '1', 10)
      if (num > max) num = max
      if (num < min) num = min
      const isExact = !value.includes('*')
      return { isExact, num }
    }

    const parsedSeconds = parseCrontabField(parts[0], 0, 59)
    const parsedMinutes = parseCrontabField(parts[1], 0, 59)
    const parsedHours = parseCrontabField(parts[2], 0, 23)
    const parsedDayOfWeek = parseCrontabField(parts[5], 0, 6)
    const parsedDayOfMonth = parseCrontabField(parts[3], 1, 31)

    setSeconds(parsedSeconds.num.toString())
    setSecondsExact(parsedSeconds.isExact)
    setMinutes(parsedMinutes.num.toString())
    setMinutesExact(parsedMinutes.isExact)
    setHours(parsedHours.num.toString())
    setHoursExact(parsedHours.isExact)
    setDayOfWeek(parsedDayOfWeek.num.toString())
    setDayOfWeekExact(parsedDayOfWeek.isExact)
    setDayOfMonth(parsedDayOfMonth.num.toString())
    setDayOfMonthExact(parsedDayOfMonth.isExact)
  }, [currentNode.metadata.config])

  const cronExpression = useMemo(() => {
    const secondExpression = () => {
      if (seconds === '*') return '*'
      if (secondsExact) return seconds
      return `*/${seconds}`
    }
    const minuteExpression = () => {
      if (minutes === '*') return '*'
      if (minutesExact) return minutes
      return `*/${minutes}`
    }
    const hourExpression = () => {
      if (hours === '*') return '*'
      if (hoursExact) return hours
      return `*/${hours}`
    }
    const dayOfMonthExpression = () => {
      if (dayOfMonth === '*') return '*'
      if (dayOfMonthExact) return dayOfMonth
      return `*/${dayOfMonth}`
    }
    const monthExpression = () => {
      if (month === '*') return '*'
      if (monthExact) return month
      return `*/${month}`
    }
    const dayOfWeekExpression = () => {
      if (dayOfWeek === '*') return '*'
      if (dayOfWeekExact) return dayOfWeek
      return `*/${dayOfWeek}`
    }
    return `${secondExpression()} ${minuteExpression()} ${hourExpression()} ${dayOfMonthExpression()} ${monthExpression()} ${dayOfWeekExpression()}`
  }, [
    seconds,
    minutes,
    hours,
    dayOfMonth,
    month,
    dayOfWeek,
    secondsExact,
    minutesExact,
    hoursExact,
    dayOfMonthExact,
    monthExact,
    dayOfWeekExact,
  ])

  // Update config and dependencies when options have been updated
  useEffect(() => {
    if (!cronExpression) return

    // Update config
    const updatedNode = structuredClone(currentNode)
    updatedNode.metadata.config = {
      cronExpression,
    }

    upsertNode(currentNode.id, updatedNode)

    // eslint-disable-next-line -- currentNode in the dependencies array causes node configs to sometimes be copied
  }, [cronExpression, upsertNode])

  useEffect(() => {
    const interval = setInterval(() => {
      setNow(new Date())
    }, 1000)

    return () => clearInterval(interval)
  }, [])

  const description = useMemo(() => {
    const dayOfWeekDescription = (() => {
      if (dayOfWeek === '*') return null
      const dayOfWeekValue = parseInt(dayOfWeek)
      if (dayOfWeekExact) {
        const dayOfWeekName = new Date(0, 0, dayOfWeekValue).toLocaleString('default', {
          weekday: 'long',
        })
        return `on ${dayOfWeekName}s`
      }
      if (dayOfWeekValue > 1) return `every ${numberToOrdinal(dayOfWeekValue)} day of the week`
      return null
    })()
    // const monthDescription = (() => {
    //   const monthName = new Date(0, monthValue - 1).toLocaleString('default', { month: 'long' })

    //   if (month === '*') return null
    //   if (monthExact) return `in ${monthName}`
    //   return `every ${monthValue} months`
    // })()
    const dayOfMonthDescription = (() => {
      if (dayOfMonth === '*') return null
      const dayOfMonthValue = parseInt(dayOfMonth)
      if (dayOfMonthExact) return `on the ${numberToOrdinal(dayOfMonthValue)} day of the month`
      if (dayOfMonthValue > 1) return `every ${numberToOrdinal(dayOfMonthValue)} day of the month`
      return null
    })()
    const hoursDescription = (() => {
      if (hours === '*') return null
      const hourValue = parseInt(hours)
      if (hoursExact) return `on the ${numberToOrdinal(hourValue)} hour`
      if (hourValue > 1) return `at every ${hourValue} hour(s)`
      return null
    })()
    const minutesDescription = (() => {
      if (minutes === '*') return null
      const minuteValue = parseInt(minutes)
      if (minutesExact) return `on the ${numberToOrdinal(minuteValue)} minute`
      if (minuteValue > 1) return `at every ${minuteValue} minute(s)`
      return null
    })()
    const secondsDescription = (() => {
      if (seconds === '*') return 'Every second'
      const secondValue = parseInt(seconds)
      if (secondsExact) return `On the ${numberToOrdinal(secondValue)} second`
      return `Every ${secondValue} second(s)`
    })()

    const fullDescription = (() => {
      let d = secondsDescription
      if (minutesDescription) d += `, ${minutesDescription}`
      if (hoursDescription) d += `, ${hoursDescription}`
      if (dayOfMonthDescription) d += `, ${dayOfMonthDescription}`
      // if (monthDescription) d += `, ${monthDescription}`
      if (dayOfWeekDescription) d += `, ${dayOfWeekDescription}`
      return d
    })()
    return fullDescription + '.'
  }, [
    seconds,
    minutes,
    hours,
    dayOfMonth,
    // month,
    dayOfWeek,
    secondsExact,
    minutesExact,
    hoursExact,
    dayOfMonthExact,
    // monthExact,
    dayOfWeekExact,
  ])

  return (
    <>
      <h2>Timer Schedule</h2>
      <p>{description}</p>
      <p>Next Runs:</p>
      <ul className="list-disc pl-5 mb-4">
        {nextExecutionTimes.map((time, index) => (
          <li key={index}>{time}</li>
        ))}
      </ul>
      <div className="flex gap-4">
        <FormFieldText
          name="seconds"
          label="Seconds"
          className="grow"
          inputType="number"
          min={'1'}
          max={'59'}
          value={seconds}
          onChange={(e) => {
            setSeconds(e.target.value)
          }}
        />
        <FormFieldSelect
          name="secondsExact"
          label="Seconds Exact"
          variant="cards"
          options={[{ label: 'Exact', value: 'exact' }]}
          optional
          hideLabel
          value={secondsExact ? 'exact' : ''}
          onChange={(e) => {
            setSecondsExact(e.target.value === 'exact')
          }}
        />
      </div>
      <div className="flex gap-4">
        <FormFieldText
          name="minutes"
          label="Minutes"
          className="grow"
          inputType="number"
          min={'1'}
          max={'59'}
          value={minutes}
          onChange={(e) => {
            setMinutes(e.target.value)
          }}
        />
        <FormFieldSelect
          name="minutesExact"
          label="Minutes Exact"
          variant="cards"
          options={[{ label: 'Exact', value: 'exact' }]}
          optional
          hideLabel
          value={minutesExact ? 'exact' : ''}
          onChange={(e) => {
            setMinutesExact(e.target.value === 'exact')
          }}
        />
      </div>
      <div className="flex gap-4">
        <FormFieldText
          name="hours"
          label="Hours"
          className="grow"
          inputType="number"
          min={'1'}
          max={'23'}
          value={hours}
          onChange={(e) => {
            setHours(e.target.value)
          }}
        />
        <FormFieldSelect
          name="hoursExact"
          label="Hours Exact"
          variant="cards"
          options={[{ label: 'Exact', value: 'exact' }]}
          optional
          hideLabel
          value={hoursExact ? 'exact' : ''}
          onChange={(e) => {
            setHoursExact(e.target.value === 'exact')
          }}
        />
      </div>
      <div className="flex gap-4">
        <FormFieldText
          name="dayOfWeek"
          label="Day of Week"
          className="grow"
          inputType="number"
          min={'1'}
          max={'7'}
          value={dayOfWeek}
          onChange={(e) => {
            setDayOfWeek(e.target.value)
          }}
        />
        <FormFieldSelect
          name="dayOfWeekExact"
          label="Day of Week Exact"
          variant="cards"
          options={[{ label: 'Exact', value: 'exact' }]}
          optional
          hideLabel
          value={dayOfWeekExact ? 'exact' : ''}
          onChange={(e) => {
            setDayOfWeekExact(e.target.value === 'exact')
          }}
        />
      </div>
      <div className="flex gap-4">
        <FormFieldText
          name="dayOfMonth"
          label="Day of Month"
          className="grow"
          inputType="number"
          min={'1'}
          max={'31'}
          value={dayOfMonth}
          onChange={(e) => {
            setDayOfMonth(e.target.value)
          }}
        />
        <FormFieldSelect
          name="dayOfMonthExact"
          label="Day of Month Exact"
          variant="cards"
          options={[{ label: 'Exact', value: 'exact' }]}
          optional
          hideLabel
          value={dayOfMonthExact ? 'exact' : ''}
          onChange={(e) => {
            setDayOfMonthExact(e.target.value === 'exact')
          }}
        />
      </div>
      {/* <div className="flex gap-4">
        <FormFieldText
          name="month"
          label="Month"
          className="grow"
          inputType="number"
          min={'1'}
          max={'12'}
          value={month}
        />
        <FormFieldSelect
          name="monthExact"
          label="Month Exact"
          variant="cards"
          options={[{ label: 'Exact', value: 'exact' }]}
          optional
          hideLabel
          value={monthExact ? 'exact' : ''}
          onChange={(e) => {
            setMonthExact(e.target.value === 'exact')
          }}
        />
      </div> */}

      <span
        className="text-gray-500 hover:text-gray-600 dark:hover:text-gray-400 cursor-pointer transition"
        onClick={(e) => {
          e.stopPropagation()
          setShowAdvanced(!showAdvanced)
        }}
      >
        Advanced
        <RemixIcon
          icon={riArrowDownSFill}
          className={`transition ml-px ${showAdvanced ? 'rotate-180' : ''}`}
        />
      </span>

      {showAdvanced && (
        <FormFieldText name="cronExpression" label="CRON Expression" value={cronExpression} />
      )}
    </>
  )
}

function numberToOrdinal(num: number) {
  const suffix =
    num % 10 === 1 && num % 100 !== 11
      ? 'st'
      : num % 10 === 2 && num % 100 !== 12
        ? 'nd'
        : num % 10 === 3 && num % 100 !== 13
          ? 'rd'
          : 'th'
  return `${num}${suffix}`
}

export default TimerSourceOptions
