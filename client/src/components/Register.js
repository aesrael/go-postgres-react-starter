import React, { useState } from 'react'
import { apiURl } from '../api'

const Register = ({ history }) => {
  const [state, setState] = useState({
    email: '',
    password: '',
    name: '',
    isSubmitting: false,
    message: '',
    errors: null,
  })

  const { email, password, name, message, isSubmitting, errors } = state

  const handleChange = async e => {
    await setState({ ...state, [e.target.name]: e.target.value })
  }

  const handleSubmit = async e => {
    setState({ ...state, isSubmitting: true })

    const { email, password, name } = state
    try {
      const res = await fetch(`${apiURl}/register`, {
        method: 'POST',
        body: JSON.stringify({
          email,
          password,
          name,
        }),
        headers: {
          'Content-Type': 'application/json',
        },
      }).then(res => res.json())
      const { success, msg, errors } = res

      if (!success) {
        return setState({ ...state, message: msg, errors, isSubmitting: false })
      }

      history.push('/login')
    } catch (e) {
      setState({ ...state, message: e.toString(), isSubmitting: false })
    }
  }

  return (
    <div className="wrapper">
      <h1>Register</h1>
      <input
        className="input"
        type="name"
        placeholder="Name"
        value={name}
        name="name"
        onChange={e => {
          handleChange(e)
        }}
      />
      <input
        className="input"
        type="text"
        placeholder="Email"
        value={email}
        name="email"
        onChange={e => {
          handleChange(e)
        }}
      />
      <input
        className="input"
        type="password"
        placeholder="Password"
        value={password}
        name="password"
        onChange={e => {
          handleChange(e)
        }}
      />

      <button disabled={isSubmitting} onClick={() => handleSubmit()}>
        {isSubmitting ? '.....' : 'Sign Up'}
      </button>
      <div className="message">{message && <p>&bull; {message}</p>}</div>
      <div>
        {errors &&
          errors.map((error, id) => {
            return <p key={id}> &bull; {error}</p>
          })}
      </div>
    </div>
  )
}

export default Register
