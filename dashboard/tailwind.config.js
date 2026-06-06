/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        background: '#0f172a',
        foreground: '#f8fafc',
        card: '#1e293b',
        'card-foreground': '#f8fafc',
        primary: '#0ea5e9',
        'primary-foreground': '#f8fafc',
        secondary: '#334155',
        'secondary-foreground': '#f8fafc',
        accent: '#38bdf8',
        'accent-foreground': '#f8fafc',
        destructive: '#ef4444',
        'destructive-foreground': '#f8fafc',
        muted: '#1e293b',
        'muted-foreground': '#94a3b8',
        border: '#334155',
        input: '#334155',
        ring: '#38bdf8',
      },
    },
  },
  plugins: [],
}
