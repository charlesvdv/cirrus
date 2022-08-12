const defaultTheme = require('tailwindcss/defaultTheme')

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      fontFamily: {
        sans: ['InterVariable', ...defaultTheme.fontFamily.sans],
      },
      colors: {
        primary: '#4E88B8',
        secondary: '#4983b3',
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms')
  ],
}
