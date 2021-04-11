import React, { FC } from 'react'
import { Box } from '@material-ui/core'
import { Field } from 'formik'
import { TextField } from 'formik-material-ui'
 
export const LoginForm: FC = () => {
  return (
    <>
      <Box my={1}>
        <Field
          component={TextField}
          name="email"
          type="email"
          label="Email"
          required
          fullWidth
        />
      </Box>

      <Box my={1}>
        <Field
          component={TextField}
          name="password"
          type="password"
          label="Password"
          required
          fullWidth
        />
      </Box>
    </>
  )
}
