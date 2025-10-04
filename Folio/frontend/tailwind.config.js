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
        display: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
      },
      colors: {
        // Dark theme inspired by Netflix/Tidal
        dark: {
          50: '#f8fafc',
          100: '#f1f5f9',
          200: '#e2e8f0',
          300: '#cbd5e1',
          400: '#94a3b8',
          500: '#64748b',
          600: '#475569',
          700: '#334155',
          800: '#1e293b',
          900: '#0f172a',
          950: '#020617',
        },
        // Accent colors
        accent: {
          red: '#e50914', // Netflix red
          blue: '#1db954', // Spotify green (Tidal-inspired)
          purple: '#8b5cf6',
          orange: '#f59e0b',
        }
      },
      spacing: {
        '18': '4.5rem',
        '88': '22rem',
      },
      animation: {
        'fade-in': 'fadeIn 0.5s ease-in-out',
        'slide-up': 'slideUp 0.3s ease-out',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { transform: 'translateY(10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
      },
    },
  },
  plugins: [require('daisyui')],
  daisyui: {
    themes: [
      {
        folio: {
          "primary": "#e50914", // Netflix red
          "secondary": "#1db954", // Spotify green
          "accent": "#8b5cf6", // Purple accent
          "neutral": "#0f172a", // Dark slate
          "base-100": "#0f172a", // Dark background
          "base-200": "#1e293b", // Slightly lighter dark
          "base-300": "#334155", // Border color
          "info": "#3b82f6",
          "success": "#10b981",
          "warning": "#f59e0b",
          "error": "#ef4444",
          "border": "#334155",
        },
      },
    ],
    base: true,
    styled: true,
    utils: true,
  },
}

