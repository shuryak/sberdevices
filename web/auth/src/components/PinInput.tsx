import React, {KeyboardEvent, useEffect, useRef, useState} from 'react'
import '@styles/PinInput.scss'

export const PinInput: React.FC<{ count: number, onPinChange: (pin: string) => void }> = ({count, onPinChange}) => {
  const [code, setCode] = useState(Array(count).fill(''))
  const inputsRef = useRef([])

  useEffect(() => {
    inputsRef.current[0].focus()
  }, [])

  useEffect(() => {
    onPinChange(code.join(''))
  }, [code])

  const onPaste = (e: React.ClipboardEvent) => {
    e.preventDefault()

    const pasteData = e.clipboardData.getData('text')

    const pasteArray = pasteData.split('').slice(0, count)
    const newCode = [...code]
    pasteArray.forEach((char, i) => {
      newCode[i] = char;
    })

    setCode(newCode)


    if (pasteArray.length > 0) {
      const nextIndex = pasteArray.length >= count ? count - 1 : pasteArray.length;
      inputsRef.current[nextIndex].focus();
    }
  }

  const onChange = (e: React.ChangeEvent<HTMLInputElement>, idx: number) => {
    const value = e.target.value
    if (value.length > 1) return

    setCode(prev => {
      const newCode = [...prev]
      newCode[idx] = value
      return newCode
    })

    if (value && idx < count - 1) {
      inputsRef.current[idx + 1].focus()
    }
  }

  const onKeyDown = (e: KeyboardEvent<HTMLInputElement>, idx: number) => {
    if (e.key === 'Backspace' && idx > 0) {
      e.preventDefault()
      setCode(prev => {
        const newCode = [...prev]
        newCode[idx] = ''
        return newCode
      })
      inputsRef.current[idx - 1].focus()
    }
    if (e.key === 'ArrowLeft' && idx > 0) {
      inputsRef.current[idx - 1].focus()
    }
    if (e.key === 'ArrowRight' && idx < count - 1) {
      inputsRef.current[idx + 1].focus()
    }
  }

  return (
    <div className="pin-input-wrapper">
      <label htmlFor="pin-input">Код из SMS:</label>
      <div className="pin-inputs" onPaste={onPaste}>
        {code.map((value, idx) => (
          <input
            key={idx}
            onChange={e => onChange(e, idx)}
            onKeyDown={e => onKeyDown(e, idx)}
            type="tel"
            maxLength={1}
            placeholder="_"
            tabIndex={1}
            inputMode="numeric"
            autoComplete={idx === 0 ? "one-time-code" : "off"}
            className="pin-input"
            ref={el => (inputsRef.current[idx] = el)}
            value={value}
          />
        ))}
      </div>
    </div>
  )
}

export default PinInput
