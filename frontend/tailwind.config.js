// Everplain Web semantic tokens. Keep upstream utility names to minimize merge churn.
const token = (name) => ({ opacityValue }) => opacityValue === undefined
  ? `var(--qx-${name})`
  : `color-mix(in srgb, var(--qx-${name}) calc(${opacityValue} * 100%), transparent)`
const colors = (names) => Object.fromEntries(Object.entries(names).map(([step, name]) => [step, token(`color-${name}`)]))
const neutral = colors({50:'surface-muted',100:'surface-strong',200:'rule',300:'rule-strong',400:'faint',500:'muted',600:'muted',700:'ink-soft',800:'ink-soft',900:'ink',950:'ink'})
const status = (name) => Object.fromEntries([50,100,200,300,400,500,600,700,800,900,950].map(step => [step, ({opacityValue}) => {
  const color = step <= 300 || step >= 900 ? `color-mix(in srgb, var(--qx-color-${name}) ${step <= 100 ? 10 : 18}%, var(--qx-color-surface))` : `var(--qx-color-${name})`
  return opacityValue === undefined ? color : `color-mix(in srgb, ${color} calc(${opacityValue} * 100%), transparent)`
}]))
/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        gray: neutral,
        slate: neutral,
        primary: colors({50:'accent-muted',100:'accent-soft',200:'rule',300:'muted',400:'accent',500:'accent',600:'accent',700:'accent-hover',800:'ink-soft',900:'accent-soft',950:'accent-muted'}),
        accent: neutral,
        dark: colors({50:'ink',100:'ink',200:'ink-soft',300:'ink-soft',400:'muted',500:'faint',600:'rule-strong',700:'rule',800:'surface',900:'surface-muted',950:'canvas'}),
        red: status('danger'), rose: status('danger'),
        green: status('success'), emerald: status('success'),
        amber: status('warning'), yellow: status('warning'), orange: status('warning'),
        blue: status('info'), cyan: status('data-teal'), sky: status('info'),
        teal: status('data-teal'), purple: status('data-violet'), indigo: status('data-blue'), violet: status('data-violet'),
      },
      fontFamily: { sans: ['var(--qx-font-ui)'], serif: ['var(--qx-font-reading)'], mono: ['var(--qx-font-mono)'] },
      fontSize: {
        xs: ['var(--qx-text-meta)', 'var(--qx-text-meta--line-height)'],
        sm: ['var(--qx-text-control)', 'var(--qx-text-control--line-height)'],
        base: ['var(--qx-text-body)', 'var(--qx-text-body--line-height)'],
        lg: ['var(--qx-text-title)', 'var(--qx-text-title--line-height)'],
        xl: ['var(--qx-text-title)', 'var(--qx-text-title--line-height)'],
        '2xl': ['var(--qx-text-section)', 'var(--qx-text-section--line-height)'],
        '3xl': ['var(--qx-text-display)', 'var(--qx-text-display--line-height)'],
      },
      boxShadow: {
        glass: 'var(--qx-shadow-panel)', 'glass-sm': 'var(--qx-shadow-card)',
        glow: 'none', 'glow-lg': 'none', card: 'var(--qx-shadow-card)',
        'card-hover': 'var(--qx-shadow-menu)', 'inner-glow': 'none',
      },
      backgroundImage: {
        'gradient-primary': 'none', 'gradient-dark': 'none',
        'gradient-glass': 'none', 'mesh-gradient': 'none',
      },
      animation: {
        'fade-in': 'fadeIn 0.3s ease-out',
        'slide-up': 'slideUp 0.3s ease-out',
        'slide-down': 'slideDown 0.3s ease-out',
        'slide-in-right': 'slideInRight 0.3s ease-out',
        'scale-in': 'scaleIn 0.2s ease-out',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        shimmer: 'shimmer 2s linear infinite',
        glow: 'none'
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' }
        },
        slideUp: {
          '0%': { opacity: '0', transform: 'translateY(10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' }
        },
        slideDown: {
          '0%': { opacity: '0', transform: 'translateY(-10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' }
        },
        slideInRight: {
          '0%': { opacity: '0', transform: 'translateX(20px)' },
          '100%': { opacity: '1', transform: 'translateX(0)' }
        },
        scaleIn: {
          '0%': { opacity: '0', transform: 'scale(0.95)' },
          '100%': { opacity: '1', transform: 'scale(1)' }
        },
        shimmer: {
          '0%': { backgroundPosition: '-200% 0' },
          '100%': { backgroundPosition: '200% 0' }
        },
        glow: {
          '0%': { boxShadow: 'var(--qx-shadow-card)' },
          '100%': { boxShadow: 'var(--qx-shadow-menu)' }
        }
      },
      backdropBlur: {
        xs: '2px'
      },
      borderRadius: {
        '4xl': 'var(--qx-radius-modal)'
      }
    }
  },
  plugins: []
}
