import React, {useEffect, useRef, useState} from 'react'

export const PhoneInput: React.FC<{
  onPhoneChange: (phone: string, isFull: boolean) => void,
  autoFocus: boolean,
  disabled: boolean
}> = ({onPhoneChange, autoFocus, disabled}) => {
  const [value, setValue] = useState<string>('+7 (')
  const [placeholder, setPlaceholder] = useState<string>('+7 (___) ___-__-__')
  const ref = useRef<HTMLInputElement>(null)

  useEffect(() => {
    autoFocus &&
    setTimeout(() => {
      if (ref.current) {
        ref.current.focus()
      }
    }, 1000)
  })

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const inputValue = e.target.value

    if (inputValue.substring(e.target.value.length - 2, e.target.value.length) == '  ') {
      return
    }

    const payload = inputValue.replace(/[+7 ()-]/g, '')

    if (!/^\d*$/.test(payload)) {
      return
    }

    if (payload.length > 10) {
      return
    } 

    let formattedValue = '+7'

    if (payload.length > 0) {
      formattedValue += ` (${payload.substring(0, 3)}`
    }

    if (payload.length >= 4) {
      formattedValue += `) ${payload.substring(3, 6)}`
    }

    if (payload.length >= 7) {
      formattedValue += `-${payload.substring(6, 8)}`
    }

    if (payload.length >= 9) {
      formattedValue += `-${payload.substring(8, 10)}`
    }

    setValue(formattedValue)
    onPhoneChange(payload, payload.length === 10)
  }

  useEffect(() => {
    const defaultPlaceholder = '+7 (___) ___-__-__'
    setPlaceholder(value + defaultPlaceholder.substring(value.length, defaultPlaceholder.length))
  }, [value])

  return (
    <>
      <label htmlFor="phone">Номер телефона:</label>
      <div className="phone-input-wrapper">
        <input
          className="phone-input"
          placeholder={placeholder}
          type="tel"
          autoComplete="off"
          disabled
        />
        <input
          id="phone"
          name="phone"
          className="phone-input"
          type="tel"
          autoComplete="off"
          value={value}
          onChange={onChange}
          ref={ref}
          disabled={disabled}
        />
      </div>
    </>
  )
}

export default PhoneInput
