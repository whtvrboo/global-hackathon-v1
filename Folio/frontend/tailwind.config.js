/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
      },
    },
  },
  plugins: [require('daisyui')],
  daisyui: {
    themes: [
      {
        folio: {
          "primary": "#5E6AD2",
          "secondary": "#26B5CE",
          "accent": "#F59E0B",
          "neutral": "#191D24",
          "base-100": "#FFFFFF",
          "base-200": "#F8FAFC",
          "base-300": "#E2E8F0",
          "info": "#3ABFF8",
          "success": "#36D399",
          "warning": "#FBBD23",
          "error": "#F87272",
          "border": "#E2E8F0",
        },
      },
    ],
    base: true,
    styled: true,
    utils: true,
  },
}

