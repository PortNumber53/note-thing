import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import tailwindcss from '@tailwindcss/vite'
import { cloudflare } from "@cloudflare/vite-plugin";

// https://vite.dev/config/
export default defineConfig({
  plugins: [tailwindcss(), react(), cloudflare()],
  resolve: {
    alias: {
      "@": "/src",
    },
  },
  server: {
    host: '0.0.0.0',
    port: 18610,
    strictPort: true,
    allowedHosts: [
      'note14.dev.portnumber53.com',
      'note16.dev.portnumber53.com',
    ],
    proxy: {
      '/api': {
        target: 'http://localhost:18611',
        changeOrigin: true,
      },
      '/auth/google': {
        target: 'http://localhost:18611',
        changeOrigin: true,
      },
      '/auth/signup': {
        target: 'http://localhost:18611',
        changeOrigin: true,
      },
      '/auth/login': {
        target: 'http://localhost:18611',
        changeOrigin: true,
      },
      '/callback': {
        target: 'http://localhost:18611',
        changeOrigin: true,
      },
    },
  },
})
