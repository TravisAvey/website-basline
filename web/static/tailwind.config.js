/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["../templates/**/*.html", "../static/src/js/editor.js"],
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["dracula", "dark", "light"],
    base: true, 
    logs: true,
    themeRoot: ":root",
    styled: true,
  },
}

