/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      /* ============================================
         CUSTOM COLOR PALETTE
         Requirements: 15.1
         Dark cinematic theme with grays, blacks, and accent colors
         ============================================ */
      colors: {
        cinema: {
          dark: {
            950: '#000000',
            900: '#0a0a0a',
            800: '#121212',
            700: '#1a1a1a',
            600: '#242424',
            500: '#2e2e2e',
            400: '#3a3a3a',
          },
          accent: {
            primary: '#e50914',
            'primary-hover': '#f40612',
            secondary: '#f5c518',
            tertiary: '#00d4ff',
          },
          text: {
            primary: '#ffffff',
            secondary: '#b3b3b3',
            tertiary: '#808080',
            muted: '#666666',
          },
        },
      },
      
      /* ============================================
         RESPONSIVE BREAKPOINTS
         Requirements: 10.1, 10.2, 10.3, 10.4
         Mobile: <768px, Tablet: 768-1023px, Desktop: 1024-1919px, Ultrawide: >=1920px
         ============================================ */
      screens: {
        'mobile': '768px',      // Mobile breakpoint (min-width)
        'tablet': '1024px',     // Tablet breakpoint
        'desktop': '1920px',    // Desktop/Ultrawide breakpoint
        'ultrawide': '2560px',  // Extra large ultrawide displays
      },
      
      /* ============================================
         CUSTOM SPACING
         Requirements: 15.2, 15.3
         Additional spacing utilities for media grids and layouts
         ============================================ */
      spacing: {
        '18': '4.5rem',   // 72px
        '22': '5.5rem',   // 88px
        '26': '6.5rem',   // 104px
        '30': '7.5rem',   // 120px
        '34': '8.5rem',   // 136px
        '38': '9.5rem',   // 152px
        '42': '10.5rem',  // 168px
        '46': '11.5rem',  // 184px
        '50': '12.5rem',  // 200px
        '128': '32rem',   // 512px
        '144': '36rem',   // 576px
      },
      
      /* ============================================
         CUSTOM SIZING
         Requirements: 10.1, 10.2, 10.3, 10.4
         Sizing utilities for responsive media cards and containers
         ============================================ */
      width: {
        'card-mobile': '150px',
        'card-tablet': '200px',
        'card-desktop': '250px',
        'card-ultrawide': '300px',
      },
      
      height: {
        'card-mobile': '225px',
        'card-tablet': '300px',
        'card-desktop': '375px',
        'card-ultrawide': '450px',
      },
      
      /* ============================================
         BACKDROP BLUR
         Requirements: 15.4
         Glassmorphism and backdrop blur utilities
         ============================================ */
      backdropBlur: {
        'cinema-sm': '4px',
        'cinema': '12px',
        'cinema-lg': '16px',
        'cinema-xl': '24px',
      },
      
      /* ============================================
         BOX SHADOWS
         Requirements: 15.2
         Cinematic shadow effects for depth
         ============================================ */
      boxShadow: {
        'cinema-sm': '0 2px 8px rgba(0, 0, 0, 0.4)',
        'cinema': '0 4px 20px rgba(0, 0, 0, 0.5)',
        'cinema-lg': '0 10px 40px rgba(0, 0, 0, 0.7)',
        'cinema-xl': '0 20px 60px rgba(0, 0, 0, 0.9)',
      },
      
      /* ============================================
         BORDER RADIUS
         Requirements: 15.2
         Consistent border radius for cards and components
         ============================================ */
      borderRadius: {
        'cinema': '0.5rem',
        'cinema-lg': '0.75rem',
        'cinema-xl': '1rem',
      },
      
      /* ============================================
         TYPOGRAPHY
         Requirements: 15.5
         Font families and sizes for cinematic interface
         ============================================ */
      fontFamily: {
        sans: ['Inter', 'system-ui', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'sans-serif'],
        display: ['Poppins', 'Inter', 'sans-serif'],
      },
      
      fontSize: {
        'hero': ['4rem', { lineHeight: '1.1', fontWeight: '700' }],
        'hero-mobile': ['2.5rem', { lineHeight: '1.1', fontWeight: '700' }],
      },
      
      /* ============================================
         TRANSITIONS
         Requirements: 15.3
         Smooth transition utilities
         ============================================ */
      transitionDuration: {
        'fast': '150ms',
        'base': '250ms',
        'slow': '350ms',
      },
      
      transitionTimingFunction: {
        'cinema': 'cubic-bezier(0.4, 0, 0.2, 1)',
      },
      
      /* ============================================
         Z-INDEX
         Layer management for overlays and modals
         ============================================ */
      zIndex: {
        'header': '40',
        'modal': '50',
        'player-controls': '45',
        'overlay': '30',
      },
      
      /* ============================================
         ASPECT RATIOS
         For media cards and video players
         ============================================ */
      aspectRatio: {
        'poster': '2/3',    // Movie poster ratio
        'backdrop': '16/9', // Widescreen ratio
        'ultrawide': '21/9', // Ultrawide ratio
      },
      
      /* ============================================
         GRID TEMPLATES
         Requirements: 10.1, 10.2, 10.3, 10.4
         Responsive grid layouts for media items
         ============================================ */
      gridTemplateColumns: {
        'media-mobile': 'repeat(auto-fill, minmax(150px, 1fr))',
        'media-tablet': 'repeat(auto-fill, minmax(200px, 1fr))',
        'media-desktop': 'repeat(auto-fill, minmax(250px, 1fr))',
        'media-ultrawide': 'repeat(auto-fill, minmax(300px, 1fr))',
      },
      
      /* ============================================
         ANIMATIONS
         Requirements: 15.3
         Custom animations for smooth interactions
         ============================================ */
      keyframes: {
        'fade-in': {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        'slide-up': {
          '0%': { opacity: '0', transform: 'translateY(20px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        'slide-down': {
          '0%': { opacity: '0', transform: 'translateY(-20px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        'scale-in': {
          '0%': { opacity: '0', transform: 'scale(0.95)' },
          '100%': { opacity: '1', transform: 'scale(1)' },
        },
        'skeleton': {
          '0%': { backgroundPosition: '200% 0' },
          '100%': { backgroundPosition: '-200% 0' },
        },
      },
      
      animation: {
        'fade-in': 'fade-in 250ms ease-in',
        'slide-up': 'slide-up 350ms ease-out',
        'slide-down': 'slide-down 350ms ease-out',
        'scale-in': 'scale-in 250ms ease-out',
        'skeleton': 'skeleton 1.5s ease-in-out infinite',
      },
      
      /* ============================================
         BACKGROUND IMAGES
         Gradient utilities for overlays
         ============================================ */
      backgroundImage: {
        'gradient-overlay': 'linear-gradient(to top, #0a0a0a 0%, rgba(10, 10, 10, 0.8) 50%, transparent 100%)',
        'gradient-overlay-bottom': 'linear-gradient(to bottom, transparent 0%, rgba(10, 10, 10, 0.8) 50%, #0a0a0a 100%)',
        'gradient-radial': 'radial-gradient(circle, var(--tw-gradient-stops))',
      },
      
      /* ============================================
         OPACITY
         Additional opacity values for overlays
         ============================================ */
      opacity: {
        '15': '0.15',
        '35': '0.35',
        '65': '0.65',
        '85': '0.85',
      },
    },
  },
  plugins: [],
}
