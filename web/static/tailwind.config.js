/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["../../app/views/*.{templ,go}"],
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["dracula", "dark"],
  },
}

