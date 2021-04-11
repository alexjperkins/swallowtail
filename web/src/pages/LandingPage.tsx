import React, { FC } from 'react'
import { Redirect } from 'react-router-dom'
import { Box, Button } from '@material-ui/core'
import styled from 'styled-components'
import { ScreenContainer } from '../components/Layout'
import { Heading, Text } from '../components/Typography/Typography'
 
import Lamp from './logo.png'

const LampImg = styled.img`
  max-width: 60%
`

export const LandingPage: FC = () => {

  const renderSignUpRedirect = () => {
    return <Redirect to='/signup' />
  }

  return(
    <ScreenContainer>
      <Box
        display="flex"
        flexDirection="column"
        minHeight={['100vh', '100%']}
        width={['100%', '500px']}
        margin="0 auto"
      >
        <Box
          display="flex"
          justifyContent="space-around"
        >
          <Box mt={6}>
            <Heading gutterBottom>{`Swallowtail`}</Heading>
            <Text gutterBottom>
              {`The Trading Bot`}
            </Text>
          </Box>

          <Box maxWidth={['50%', '200px']}>
            <LampImg src={Lamp} alt="logo" />
          </Box>
        </Box>

        <Box mt={6}>
          <Text gutterBottom>
            {`Coming Soon...`}
          </Text>
        </Box>

        <Box mt={4}>
          <Box>
            <Button
              disableElevation
              variant="contained"
              color="secondary"
              onClick={renderSignUpRedirect}
              fullWidth
            >
              {`Signup`}
            </Button>
          </Box>
        </Box>
      </Box>
    </ScreenContainer>
  )
}
