import React, { useState } from 'react'
import { apiURl } from '../api'

const createCookie = (cookieName, cookieValue, hourToExpire) => {
  const date = new Date()
  date.setTime(date.getTime() + hourToExpire * 60 * 60 * 1000)
  document.cookie = `${cookieName}' = '${cookieValue}'; expires = ' ${date.toGMTString()}`
}

const Login = ({ history }) => {
  const [state, setState] = useState({
    email: '',
    password: '',
    isSubmitting: false,
    message: '',
  })

  const { email, password, isSubmitting, message } = state

  const handleChange = async e => {
    const { name, value } = e.target
    await setState({ ...state, [name]: value })
  }

  const handleSubmit = async () => {
    setState({ ...state, isSubmitting: true })
    console.log(isSubmitting)
    const { email, password } = state
    try {
      const res = await fetch(`${apiURl}/login`, {
        method: 'POST',
        body: JSON.stringify({
          email,
          password,
        }),
        headers: {
          'Content-Type': 'application/json',
        },
      }).then(res => res.json())

      const { token, success, msg, user } = res

      if (!success) {
        return setState({
          ...state,
          message: msg,
          isSubmitting: false,
        })
      }
      createCookie('token', token, 0.4)
      history.push({ pathname: '/session', state: user })
    } catch (e) {
      setState({ ...state, message: e.toString(), isSubmitting: false })
    }
  }

  return (
    <div className="wrapper">
      <h1>Login</h1>
      <input
        className="input"
        type="text"
        placeholder="email"
        value={email}
        name="email"
        onChange={e => {
          handleChange(e)
        }}
      />

      <input
        className="input"
        type="password"
        placeholder="password"
        value={password}
        name="password"
        onChange={e => {
          handleChange(e)
        }}
      />

      <button disabled={isSubmitting} onClick={() => handleSubmit()}>
        {isSubmitting ? '.....' : 'login'}
      </button>
      <div className="message">{message}</div>
    </div>
  )
}

export default Login
