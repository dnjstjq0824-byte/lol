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
          dark: '#0B0F19',
          card: '#111827',
          border: '#1F2937',
          accent: '#06B6D4',
          victory: '#3B82F6',
          defeat: '#EF4444',
          text: '#F9FAFB',
          muted: '#9CA3AF'
        }
      }
    },
  },
  plugins: [],
}

