import { ApolloClient } from 'apollo-client'
import { InMemoryCache } from 'apollo-cache-inmemory'
import { HttpLink } from 'apollo-link-http'
import { onError } from 'apollo-link-error'
import { ApolloLink } from 'apollo-link'

const DOMAIN = '0.0.0.0:5000'
const ENDPOINT = '/graphql'
const GRAPHQL_ENDPOINT = `${DOMAIN}${ENDPOINT}`

const errorLink = onError(props => {
	const { graphQLErrors, networkError } = props
	if (graphQLErrors)
		graphQLErrors.forEach(({ message, locations, path }) =>
			console.log(
				`[GraphQL error]: Message: ${message}, Location: ${locations}, Path: ${path}`
			)
		)
     if (networkError) console.log(`[Network error]: ${networkError}`)
})


export const gqlClient = new ApolloClient({
  link: ApolloLink.from([
    errorLink,
    new HttpLink({
      uri: GRAPHQL_ENDPOINT,
      credentials: 'omit',
    })
  ]),
  cache: new InMemoryCache()
})
