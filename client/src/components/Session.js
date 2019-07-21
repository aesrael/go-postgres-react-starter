import React, { useState, useEffect } from 'react'
import { apiURl } from '../api'


const Login = ({ history }) => {
  const [state, setState] = useState({
    isFetching: false,
    message: null,
    user: null,
  })

  const { isFetching, message, user = {} } = state

  const getUserInfo = async () => {
    setState({ ...state, isFetching: true, message: 'fetching details...' })
    try {
      const res = await fetch(`${apiURl}/session`, {
        method: 'GET',
        credentials:"same-origin",
        headers: {
          Accept: 'application/json',
          Authorization: document.cookie,
        },
      }).then(res => res.json())

      const { success, user } = res
      if (!success) {
        history.push('/login')
      }
      setState({ ...state, user, message: null, isFetching: false })
    } catch (e) {
      setState({ ...state, message: e.toString(), isFetching: false })
    }
  }

  useEffect(() => {
    if (history.location.state) {
      return setState({ ...state, user: { ...history.location.state } })
    }
    getUserInfo()
  }, [])

  return (
    <div className="wrapper">
      <h1>Welcome {user && user.name}</h1>
      {user && user.email}
      <div className="message">
        {isFetching ? 'fetching details..' : message}
      </div>

      <button
        onClick={() => {
          localStorage.clear()
          history.push('/login')
        }}
      >
        Logout?
      </button>
    </div>
  )
}

export default Login
