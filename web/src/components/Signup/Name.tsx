import React, { FC } from 'react'
import{ Box } from '@material-ui/core'
import { Field } from 'formik'
import { TextField } from 'formik-material-ui'
import { Heading, Text } from '../Typography/Typography'

interface INameProps {
  heading: string
  text: string
}

export const Name: FC<INameProps> = ({ heading, text }) => {
  return (
    <>
      <Heading gutterBottom>{heading}</Heading>

      <Text gutterBottom>{text}</Text>

      <Box my={1}>
        <Field
          component={TextField}
          name="firstName"
          type="text"
          label="First name"
          required
          fullWidth
        />
      </Box>

      <Box my={1}>
        <Field
          component={TextField}
          name="lastName"
          type="text"
          label="Last name"
          required
          fullWidth
        />
      </Box>
    </>
  )
}
