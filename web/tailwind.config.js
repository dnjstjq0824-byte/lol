/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        rift: {
          dark: '#070B14',
          card: '#0D1424',
          raised: '#131C31',
          border: '#1E2A45',
          gold: '#C8AA6E',
          goldlight: '#E9D8A6',
          golddeep: '#5C4A26',
          victory: '#3B82F6',
          defeat: '#EF4444',
          text: '#F4F1EA',
          muted: '#9AA3B2',
        }
      },
      fontFamily: {
        display: ['"Segoe UI"', 'system-ui', '-apple-system', 'Roboto', '"Noto Sans KR"', 'sans-serif'],
      },
      boxShadow: {
        card: '0 1px 2px rgba(0, 0, 0, 0.5), 0 8px 24px -12px rgba(0, 0, 0, 0.6)',
        pop: '0 12px 32px -8px rgba(0, 0, 0, 0.7)',
      },
    },
  },
  plugins: [],
}
