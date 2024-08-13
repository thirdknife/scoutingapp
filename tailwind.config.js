/** @type {import('tailwindcss').Config} */
module.exports = {
  mode: 'jit',
  content: ['./views/*.templ'],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms')
  ],
}

