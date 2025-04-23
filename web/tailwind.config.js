/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{vue,ts}"],
  theme: {
   
    dropShadow: {
      'small': '0 35px 35px rgba(0, 0, 0, 0.25)',
      'large': [
          '4px 6px 1px rgba(0, 0, 0, 0.25)'
      ]
    },
    boxShadow:{
      'large': '0px 4px 4px 0px rgba(0, 0, 0, 0.25)'
    },
    extend: {
      colors: {
        'main': '#627bd9',
        'main-light': '#6BAEFF',
        'bg': '#F9F6F9',
        'playing': '#6AB760',
      },
      borderRadius:{
        DEFAULT: '10px',
        'large': '12px'
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}