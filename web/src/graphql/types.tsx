import gql from 'graphql-tag';
import * as ApolloReactHooks from '@apollo/react-hooks';
export type Maybe<T> = T | null
 
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
  DateTime: any;
  UUID: any;
  GenericScalar: any;
}

export type Mutations = {
  __typename?: 'Mutations';
  registerUser?: Maybe<RegisterUser>;
}

export type RegisterUser = {
  __typename?: "RegisterUser";
  userId?: Maybe<Scalars['String']>;
  token?: Maybe<Scalars['String']>;
}

export type RegisterUserMutation = (
  { __typename?: 'Mutations' }
  & { registerUser?: Maybe<(
    { __typename?: 'RegisterUser' }
    & Pick<RegisterUser, 'token'>
  )> }
);
 
export type RegisterUserMutationVariables = {
  firstName: Scalars['String'],
  lastName: Scalars['String'],
  email: Scalars['String'],
  password: Scalars['String'],
}

export const RegisterUserDocument = gql`
    mutation registerUser(
      $firstName: String!,
      $lastName: String!,
      $email: String!,
      $password: String!,
    ) {
  registerUser(
      firstName: $firstName,
      lastName: $lastName,
      email: $email,
      password: $password,
      ) {
      token
  }
}
`

/**
 * __useRegisterUserMutation__
 *
 * To run a mutation, you first call `useRegisterUserMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRegisterUserMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [registerUserMutation, { data, loading, error }] = useRegisterUserMutation({
 *   variables: {
 *      firstName: // value for 'firstName'
 *      lastName: // value for 'lastName'
 *      email: // value for 'email'
 *      password: // value for 'password'
 *   },
 * });
 */
export function useRegisterUserMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<RegisterUserMutation, RegisterUserMutationVariables>) {
        return ApolloReactHooks.useMutation<RegisterUserMutation, RegisterUserMutationVariables>(RegisterUserDocument, baseOptions);
      }
