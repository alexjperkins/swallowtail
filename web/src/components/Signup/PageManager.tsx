import React from 'react'
import { MemoryRouter, Route, Redirect, Switch } from 'react-router-dom'
import { NavigationButtons } from './NavigationButtons'
import { Box } from '@material-ui/core'
import { ScreenContainer } from '../Layout'

import { randomPath } from './util'


interface IPage {
  Component: React.ReactNode
  fields: string[]
}

interface IPageWithPath extends IPage {
  path: string
}

export class PageManager{
  pages: IPageWithPath[]
  goBack: () =>  void

  constructor(pages: IPage[], goBack: () => void) {
    this.pages = pages.map((page: IPage) => ({...page, path: randomPath() }))
    this.goBack = goBack
  }

  get firstPageLink(): string {
    return this.pages[0].path
  }

  isFirstPage(currentPageLink: string): boolean {
    return currentPageLink === this.firstPageLink
  }

  currentPageFields(currentPageLink: string): string[] {
    const currentPage = this.pages.find(
      (page: IPageWithPath) => page.path === currentPageLink
    )

    return currentPage?.fields || []
  }

  nextPageLink(currentPageLink: string): string | null {
    const currentPageIndex = this.pages.findIndex(
      (page: IPageWithPath) => page.path === currentPageLink
    )

    const nextPage = this.pages[currentPageIndex + 1]

    if (!nextPage) {
      return null
    }

    return nextPage.path
  }

  get children(): React.ReactNode {
    return (
      <ScreenContainer>
        <Box
          justifyContent="space-between"
          display="flex"
          flexDirection="column"
          minHeight={['100vh', '100%']}
          width={['100%', '500px']}
          margin="0 auto"
        >

          <MemoryRouter>
            <Box mt={8}>
              <Switch>
                {this.pages.map((page: IPageWithPath) => (
                  <Route path={page.path} key={page.path}>
                    {page.Component}
                  </Route>
                ))}
                <Redirect to={this.firstPageLink}/>
              </Switch>
            </Box>

            <Box>
              <NavigationButtons goBack={this.goBack} pages={this} />
            </Box>

          </MemoryRouter>
        </Box>
      </ScreenContainer>
    )
  }
}
