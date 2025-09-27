/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        'fluent-purple': '#A020F0',
        'fluent-purple-light': '#B24CF0',
        'fluent-purple-dark': '#8010D0',
        'fluent-bg-light': '#F9F9F9',
        'fluent-bg-dark': '#202020',
        'fluent-surface-light': '#FFFFFF',
        'fluent-surface-dark': '#2D2D2D',
        'fluent-text-light': '#0F0F0F',
        'fluent-text-dark': '#E5E5E5',
        'fluent-border-light': '#E5E5E5',
        'fluent-border-dark': '#404040'
      },
      boxShadow: {
        'fluent-sm': '0 2px 4px rgba(0, 0, 0, 0.05)',
        'fluent': '0 4px 8px rgba(0, 0, 0, 0.08)',
        'fluent-lg': '0 8px 16px rgba(0, 0, 0, 0.12)',
      },
      transitionProperty: {
        'width': 'width',
        'spacing': 'margin, padding',
      },
      transitionTimingFunction: {
        'fluent': 'cubic-bezier(0.1, 0.9, 0.2, 1)',
      },
      transitionDuration: {
        '250': '250ms',
      }
    },
  },
  plugins: [],
}