import React, { useState, useEffect } from 'react'
import { apiURl } from '../api'
import { deleteCookie } from '../utils'

const Session = ({ history }) => {
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
        credentials: 'include',
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

  const handleLogout = () => {
    deleteCookie('token')
    history.push('/login')
  }

  useEffect(() => {
    if (history.location.state) {
      return setState({ ...state, user: { ...history.location.state } })
    }
    getUserInfo()
  }, [])

  return (
    <div className="wrapper">
      <h1>Welcome, {user && user.name}</h1>
      {user && user.email}
      <div className="message">
        {isFetching ? 'fetching details..' : message}
      </div>

      <button
        style={{ height: '30px' }}
        onClick={() => {
         handleLogout()
        }}
      >
        logout
      </button>
    </div>
  )
}

export default Session
