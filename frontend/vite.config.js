import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
// https://vitejs.dev/config/
export default defineConfig({
    plugins: [react()],
    server: {
        host: true,
        port: 28519,
        proxy: {
            '/api': {
                target: 'http://localhost:29519',
                changeOrigin: true,
            },
            '/ws': {
                target: 'ws://localhost:29519',
                ws: true,
            },
            '/healthz': {
                target: 'http://localhost:29519',
                changeOrigin: true,
            },
        },
    },
});
