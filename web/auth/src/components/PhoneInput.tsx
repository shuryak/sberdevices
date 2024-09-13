import React, { useEffect, useRef, useState} from 'react'
import InputMask from 'react-input-mask'

export const PhoneInput: React.FC<{
  onPhoneChange: (phone: string, isFull: boolean) => void,
  autoFocus: boolean,
  disabled: boolean
}> = ({onPhoneChange, autoFocus, disabled}) => {
  const [value, setValue] = useState<string>('')
  const [placeholder, setPlaceholder] = useState<string>('+7 (___) ___-__-__')
  const ref = useRef<HTMLInputElement>(null)

  useEffect(() => {
    autoFocus &&
    setTimeout(() => {
      if (ref.current) {
        ref.current.focus()
        ref.current.click()
      }
    }, 1000)
  })

  useEffect(() => {
    const defaultPlaceholder = '+7 (___) ___-__-__'

    setPlaceholder(value + defaultPlaceholder.substring(value.length, defaultPlaceholder.length))
  }, [value])

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const payload = e.target.value.replace(/[+7 _()-]/g, '')

    setValue(e.target.value)
    onPhoneChange(payload, payload.length == 10)
  }

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
        <InputMask
          mask="+7 (999) 999-99-99"
          value={value}
          onChange={onChange}
          alwaysShowMask={true}
          maskPlaceholder={null}
          id="phone"
          name="phone"
          className="phone-input"
          type="tel"
          autoComplete="off"
          autoFocus={autoFocus}
          inputRef={ref}
          disabled={disabled}
        />
      </div>
    </>
  )
}

export default PhoneInput
