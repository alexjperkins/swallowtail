import React, { FC } from 'react'
import { Typography, TypographyProps } from '@material-ui/core'
import styled from 'styled-components'

export const Heading: FC<TypographyProps> = props => (
  <Typography color="primary" variant="h4" {...props} />
)

const HeadingSerifTypography = styled(Typography)({
  fontFamily: 'PT Serif'
})

export const HeadingSerif: FC<TypographyProps> = props => (
  <HeadingSerifTypography color="primary" variant="subtitle1" {...props} />
)

export const HeadingSmall: FC<TypographyProps> = props => (
  <Typography color="primary" variant="h6" {...props} />
)

export const Text: FC<TypographyProps> = props => (
  <Typography color="textPrimary" variant="body2" {...props} />
)

export const TextNormal: FC<TypographyProps> = props => (
  <Typography color="textPrimary" variant="body1" {...props} />
)
