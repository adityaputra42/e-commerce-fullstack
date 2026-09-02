import type { Config } from "tailwindcss"

const config = {
  content: [
    './pages/**/*.{ts,tsx}',
    './components/**/*.{ts,tsx}',
    './app/**/*.{ts,tsx}',
    './src/**/*.{ts,tsx}',
  ],
  theme: {
    container: {
      center: true,
      padding: "2rem",
      screens: {
        "2xl": "1400px",
      },
    },
    extend: {
      colors: {
        primary: {
          DEFAULT: '#DC95FF',
          dark: '#3E0F8D',
          light: '#7692FF',
          foreground: '#ffffff',
        },
        secondary: {
          DEFAULT: '#18181b',
          dark: '#000000',
          light: '#3f3f46',
          foreground: '#ffffff',
        },
        neutral: {
          50: '#faf9f7',
          100: '#f2f1ee',
          200: '#e5e3de',
          300: '#d4d1c9',
          400: '#a3a099',
          500: '#78766f',
          600: '#57554f',
          700: '#3f3d38',
          800: '#29281f',
          900: '#18181b',
        },
        success: '#16a34a',
        warning: '#d97706',
        error: '#dc2626',
        info: '#2563eb',
      },
      fontFamily: {
        sans: ['Inter', 'Poppins', 'sans-serif'],
        display: ['Poppins', 'sans-serif'],
      },
      borderRadius: {
        'xl': '0.625rem',
        '2xl': '0.75rem',
        '3xl': '1rem',
      },
    },
  },
  plugins: [],
} satisfies Config

export default config
