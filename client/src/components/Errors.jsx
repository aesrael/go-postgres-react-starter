import React from 'react'

export default ({ errors }) => (
  <div>
    {errors.map((error) => (
      <li key={error} className="errors">
        {error}
      </li>
    ))}
  </div>
)
