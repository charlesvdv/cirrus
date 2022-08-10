import { defineConfig } from 'vite';
import solidPlugin from 'vite-plugin-solid';

export default defineConfig({
  plugins: [solidPlugin()],
  server: {
    port: 3000,
    // Required if vite is behind HTTP proxy
    // See https://vitejs.dev/config/server-options.html#server-hmr for more info.
    hmr: { port: 3000 },
  },
  build: {
    target: 'esnext',
  },
});
