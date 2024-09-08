import PhoneInput from '@components/PhoneInput.tsx'
import PinInput from '@components/PinInput.tsx'
import '@styles/Card.scss'
import Logo from '@assets/logo.png'
import React, { useEffect, useState } from 'react'
import PulseLoader from 'react-spinners/PulseLoader'
import axios from 'axios'

interface Phone {
  phone: string;
  ok: boolean;
}

interface OAuthStartResponse {
  code: string;
}

interface OAuthOTPResponse {
  ok: boolean;
}

const BASE_URL = import.meta.env.VITE_API_BASE_URL || 'https://sdprovider.ru'

export const Card = () => {
  const [smsSent, setSmsSent] = useState<boolean>(false)
  const [phone, setPhone] = useState<Phone>({
    phone: '',
    ok: false
  })
  const [oauthCode, setOAuthCode] = useState<string>(null)
  const [otp, setOTP] = useState<string>('')

  useEffect(() => {
    if (otp.length != 5 || oauthCode === null) {
      return
    }

    axios.get<OAuthOTPResponse>(`${BASE_URL}/api/v1.0/oauth/otp?otp=${otp}&code=${oauthCode}`)
      .then(({ data }) => {
        if (!data.ok) {
          return setOTP('')
        }

        const searchParams = new URLSearchParams(window.location.search)
        const redirect = new URL(searchParams.get('redirect_uri'))
        const state = searchParams.get('state')

        const redirectSearchParams = new URLSearchParams()
        redirectSearchParams.set('state', state)
        redirectSearchParams.set('code', oauthCode)

        redirect.search = redirectSearchParams.toString()

        window.location.replace(redirect)
      })
  }, [otp])

  const sendSms = (e: React.FormEvent) => {
    e.preventDefault()

    setSmsSent(true)

    axios.get<OAuthStartResponse>(`${BASE_URL}/api/v1.0/oauth/start?phone=7${phone.phone}`)
      .then(({ data }) => {
        setOAuthCode(data.code)
      })
  }

  return (
    <div className="card">
      <img className="card-logo" src={Logo} alt="Логотип SberDevices" />

      <form onSubmit={sendSms}>
        {!oauthCode ?
          <PhoneInput autoFocus disabled={smsSent} onPhoneChange={(phone, ok) => {
            setPhone({ phone: phone, ok: ok })
          }} />
          :
          <PinInput count={5} onPinChange={setOTP} />
        }

        {!oauthCode &&
          <button type="submit" disabled={!phone.ok}>
            {smsSent ? <PulseLoader color="white" size={5} /> : 'Отправить SMS-код'}
          </button>
        }
      </form>
    </div>
  )
}
