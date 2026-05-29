/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./html/**/*.go", "./http/**/*.go"],
  theme: {
    extend: {
      colors: {
        primary: {
          300: '#c4b5fd',
          400: '#a78bfa',
          500: '#8b5cf6',
        },
        seafoam: {
          400: '#50d9a8',
          500: '#3db88a',
        },
        canary: {
          300: '#fde047',
          400: '#facc15',
          500: '#eab308',
        },
      },
    },
  },
  plugins: [],
}
