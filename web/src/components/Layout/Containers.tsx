import React, { FC } from 'react'
import { Container } from '@material-ui/core'

interface IScreenContainerProps {
  children: React.ReactNode
}

export const ScreenContainer: FC<IScreenContainerProps> = ({ children }) => {
  return <Container maxWidth="lg">{children}</Container>
}
