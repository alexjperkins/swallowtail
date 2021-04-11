import { createMuiTheme } from '@material-ui/core/styles'

export const theme = createMuiTheme({
  palette: {
    primary: {
      main: '#159CC2',
      contrastText: '#ffffff',
    },
    secondary: {
      main: '#F719FF',
      contrastText: '#ffffff',
      dark: '#000000',
    },
    text: {
      primary: 'rgba(0, 0, 0, 0.6)',
      secondary: '#FFFFFF'
    },
    error: {
      main: '#f3794e',
    },
    background: {
      default: '#7686b0',
    }
  },
  props: {
    MuiTypography: {
      variantMapping: {
        subtitle1: 'span',
      }
    }
  },
  overrides: {
    MuiFormControlLabel: {
      label: {
        color: 'rgba(0, 0, 0, 0.6)',
      }
    }
  },
})
