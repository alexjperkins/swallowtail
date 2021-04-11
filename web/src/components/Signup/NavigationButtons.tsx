import React, { FC } from 'react'
import { Box, Button, CircularProgress } from '@material-ui/core'
import { useHistory, useLocation } from 'react-router-dom'
import { useFormikContext } from 'formik'

import { PageManager } from './PageManager'

interface INavigationButtons {
  goBack: () => void
  pages: PageManager
}

export const NavigationButtons: FC<INavigationButtons> = ({
  goBack,
  pages,
}) => {

  const history = useHistory()
  const location = useLocation()
  const {
    submitForm,
    values,
    validateForm,
    resetForm,
    isSubmitting,
  } = useFormikContext()

  const nextPath = React.useMemo(() => pages.nextPageLink(location.pathname), [
    pages,
    location,
  ])

  const goToNextPage = React.useCallback(async () => {
    const errors = await validateForm()
    const currentPageFields = pages.currentPageFields(location.pathname)
    const isNotValid = currentPageFields.some(
      (fieldName: string) => (errors as any)[fieldName]
    )

    if (isNotValid) {
      submitForm()
      return
    }

    if (nextPath) {
      resetForm({ values: values as any })
      return history.push(nextPath)
    }

    return submitForm()
  }, [
    nextPath,
    location,
    pages,
    values,
    validateForm,
    submitForm,
    history,
    resetForm
  ])

  const stepBack = React.useCallback(() => {
    const isFirstPage = pages.isFirstPage(location.pathname)
    if (isFirstPage){
      goBack()
    }

    history.goBack()
  }, [history, location, pages, goBack])

  const buttonText = React.useMemo(() => {
    if (isSubmitting) {
      return <CircularProgress size={24} />
    }
    return nextPath ? 'Continue': 'Finish Registration'
  }, [isSubmitting, nextPath])


  return (
    <>
      <Box>
        <Button
          disableElevation
          variant="contained"
          color="secondary"
          onClick={goToNextPage}
          disabled={isSubmitting}
          fullWidth
        >
          {buttonText}
        </Button>
      </Box>

      <Box mt={1}>
        <Button
          disableElevation
          variant="outlined"
          color="secondary"
          onClick={stepBack}
          disabled={isSubmitting}
          fullWidth
        >
          Back
        </Button>
      </Box>
    </>
  )
}
