import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import svgr from "vite-plugin-svgr";
import viteTsconfigPaths from 'vite-tsconfig-paths'


export default defineConfig({
    plugins: [react(), svgr()],
    build: {
      outDir: "build",
    },
    server: {
        proxy: {
            '/api': {
                target: 'http://localhost:8000',
                changeOrigin: true,
                secure: false
            }
        }
    }
  });
