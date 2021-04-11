import React from 'react'
import styled from 'styled-components'
import { Box, Button } from '@material-ui/core'
import { ScreenContainer } from '../Layout'


const NavbarButton = styled(Button)`
  margin: 0;
  padding: 16;
  line-height: 1rem;
  min-width: 0;
  font-size: 0.7em;
  font-weight: bold;
`

export const Navbar = () => {
  return (
    <ScreenContainer>
      <Box
        display="flex"
        justifyContent="space-between"
        flexDirection="row"
      >

        <Box
          display="flex"
          justifyContent="flex-start"
        >
          <Box pt={1}>
            <NavbarButton>
              Home
            </NavbarButton>
          </Box>
        </Box>

        <Box
          display="flex"
          justifyContent="flex-end"
        >
          <Box pt={1}>
            <NavbarButton>
              Login
            </NavbarButton>
            <NavbarButton>
              Register
            </NavbarButton>
          </Box>
        </Box>

      </Box>
    </ScreenContainer>
  )
}
