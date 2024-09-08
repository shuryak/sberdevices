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
        ref.current.focus();
      }
    }, 1000);
  });

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.value.substring(e.target.value.length - 2, e.target.value.length) == '  ') {
      return
    }

    const payload = e.target.value.replace(/[+7 ()-]/g, '')

    if (!/^\d*$/.test(payload)) {
      return
    }

    // TODO: refactor :)
    if (payload.length > 10) {
      return
    } else if (payload.length < 3) {
      setValue(`+7 (${payload}`)
    } else if (payload.length === 3 && value[value.length - 1] == ' ') {
      setValue(`+7 (${payload.slice(0, payload.length - 1)}`)
    } else if (payload.length === 3) {
      setValue(`+7 (${payload}) `)
    } else if (payload.length > 3 && payload.length < 6) {
      setValue(`+7 (${payload.substring(0, 3)}) ${payload.substring(3, payload.length)}`)
    } else if (payload.length === 6 && value[value.length - 1] == '-') {
      setValue(`+7 (${payload.substring(0, 3)}) ${payload.substring(3, payload.length - 1)}`)
    } else if (payload.length === 6) {
      setValue(`+7 (${payload.substring(0, 3)}) ${payload.substring(3, payload.length)}-`)
    } else if (payload.length === 8 && value[value.length - 1] == '-') {
      setValue(`+7 (${payload.substring(0, 3)}) ${payload.substring(3, 6)}-${payload.substring(6, payload.length - 1)}`)
    } else if (payload.length === 8) {
      setValue(`+7 (${payload.substring(0, 3)}) ${payload.substring(3, 6)}-${payload.substring(6, 8)}-`)
    } else if (payload.length >= 9) {
      setValue(`+7 (${payload.substring(0, 3)}) ${payload.substring(3, 6)}-${payload.substring(6, 8)}-${payload.substring(8, payload.length)}`)
    } else if (payload.length >= 7) {
      setValue(`+7 (${payload.substring(0, 3)}) ${payload.substring(3, 6)}-${payload.substring(6, payload.length)}`)
    }

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
