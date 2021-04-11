import React, { FC } from 'react'
import { LoginForm } from '../LoginForm'
import { Heading, Text } from '../Typography/Typography'

export const EmailAndPassword: FC = () => {
  return (
    <>
      <Heading gutterBottom>Signup</Heading>
      <Text gutterBottom>Please fill out details here...</Text>

      <LoginForm />
    </>
  )
}
