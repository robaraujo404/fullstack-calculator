import react from '@vitejs/plugin-react'
import { defineConfig } from 'vitest/config'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    // Same-origin in development too, so the API needs no CORS.
    proxy: { '/api': 'http://localhost:8080' },
  },
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
    coverage: {
      provider: 'v8',
      // Without an explicit include, files that no test imports are left out
      // of the report and the percentage looks better than it is.
      include: ['src/**/*.{ts,tsx}'],
      exclude: ['src/**/*.test.{ts,tsx}', 'src/test/**', 'src/main.tsx'],
      // The second text reporter writes coverage/coverage.txt (every file,
      // including fully covered ones), which `make coverage` copies to docs/.
      reporter: ['text', 'html', 'json-summary', ['text', { file: 'coverage.txt', skipFull: false }]],
    },
  },
})
