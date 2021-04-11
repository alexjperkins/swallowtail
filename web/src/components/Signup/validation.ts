import * as Yup from 'yup'
import { RegisterUserMutationVariables } from  '../../graphql/types'

export const requiredString = (message: string) => 
  Yup.string()
    .trim()
    .required(message)

export const parseValuesForSignUpForm = (
  values: any
): RegisterUserMutationVariables => {
  return {
    firstName: values.firstName,
    lastName: values.lastName,
    email: values.email,
    password: values.password,
  }
}
