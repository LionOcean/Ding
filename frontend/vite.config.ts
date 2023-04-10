import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { resolve } from 'path';
import { readFileSync } from 'fs';

let wailsJSON = readFileSync(resolve(__dirname, '../wails.json'), 'utf-8');
const __APP_VERSION__ = JSON.parse(wailsJSON).info.productVersion;

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      '@utils': resolve(__dirname, 'src/utils'),
      '@wailsjs': resolve(__dirname, 'wailsjs'),
    },
  },
  define: {
    __APP_VERSION__: JSON.stringify(__APP_VERSION__),
  }
});
