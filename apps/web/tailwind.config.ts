import type { Config } from 'tailwindcss'

const config: Config = {
  content: [
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    './features/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        studio: {
          bg: '#090D16',
          card: '#111726',
          panel: '#1A2335',
          border: '#26334D',
          accent: '#06B6D4',
          violet: '#8B5CF6',
          pink: '#EC4899',
          amber: '#F59E0B',
          green: '#10B981',
          text: '#F1F5F9',
          muted: '#94A3B8',
        }
      }
    },
  },
  plugins: [],
}
export default config
