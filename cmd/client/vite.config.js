import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      '/events': {
        target: 'http://localhost:8089',
        changeOrigin: true,
      },
      '/metrics': {
        target: 'http://localhost:8089',
        changeOrigin: true,
      }
    }
  }
})

