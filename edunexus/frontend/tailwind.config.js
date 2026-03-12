/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        cyber: {
          bg: '#090D14',
          cyan: '#00FFCC',
          purple: '#B026FF',
          dark: '#1A202C',
          panel: 'rgba(26, 32, 44, 0.7)'
        }
      },
      fontFamily: {
        mono: ['"Fira Code"', 'monospace'],
        sans: ['"Inter"', 'sans-serif']
      }
    },
  },
  plugins: [],
}
