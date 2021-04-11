import React, { FC } from 'react'
import { useHistory } from 'react-router-dom'
import { Formik } from 'formik'
import * as Yup from 'yup'
import gql from 'graphql-tag'

import { useRegisterUserMutation } from '../graphql/types'

import {
  Name,
  EmailAndPassword,
  PageManager,
  parseValuesForSignUpForm,
  requiredString, 
} from '../components/Signup'
import { setAccessToken } from '../auth'


const PAGES = [
  {
    Component: (
      <Name 
        text="
        Please enter your name"
        heading="
        Personal Details"
      />
    ),
    fields: ['firstName', 'lastName']
  },
  {
    Component: (
      <EmailAndPassword />
    ),
    fields: ['email', 'password']
  }
]

export const registerUser = gql`
  mutation registerUser(
    $firstName: String!
    $lastName: String!
    $email: String!
    $password: String!
) {
  registerUser(
    firstName: $firstName
    lastName: $lastName
    email: $email
    password: $password
  ) {
    token
  }
}
`

export const SignupForm: FC = () => {

  const [ registerUserMutation ] = useRegisterUserMutation()
  const history = useHistory()

  const goBack = React.useCallback(() => history.push('/'), [ history ])
  const pages = React.useMemo(() => new PageManager(PAGES, goBack), [ goBack ])

  return (
    <Formik
      onSubmit={async (values, helpers) => {
        const variables = parseValuesForSignUpForm(values)
        try {
          const { data } = await registerUserMutation({ variables })
          helpers.setSubmitting(false)
          setAccessToken(data?.registerUser?.token!)
          history.push('/')
        } catch (error) {
          helpers.setSubmitting(false)
          console.log(error)
        }
      }}
      initialValues={{
        firstName: '',
        lastName: '',
        email: '',
        password: '',
      }}
      isInitialValid={false}
      validationSchema={Yup.object({
        firstName: requiredString('First name is required'),
        lastName: requiredString('Last name is required'),
        email: Yup.string()
          .trim()
          .required('Email is required')
          .email('Must be a valid email'),
        password: Yup.string()
          .required('Please provide a password')
          .min(8, 'Password is too short - it should be at least 8 characters')
          .matches(/[a-zA-Z]/, 'Password can only contain Latin Characters')
      })}
    >
      {pages.children}
    </Formik>
  )
}
