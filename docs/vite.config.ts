import { defineConfig } from 'vite';
import { resolve } from 'path'
import { createHtmlPlugin } from 'vite-plugin-html'

import { dirname } from 'path'
import { fileURLToPath } from 'url'

var __dirname;

const _dirname = typeof (__dirname) !== 'undefined'
  ? __dirname
  : dirname(fileURLToPath(import.meta.url))

export default defineConfig({
  plugins: [
    createHtmlPlugin({
    minify: {
      collapseWhitespace: true,
      keepClosingSlash: true,
      removeComments: true,
      removeRedundantAttributes: true,
      removeScriptTypeAttributes: false,
      removeStyleLinkTypeAttributes: false,
      useShortDoctype: true,
      minifyCSS: true,
    },
  })],
  server: {
    port: 3000,
    
  },
  build: {
    target: 'esnext',
    rollupOptions: {
      input: {
        main: resolve(_dirname, 'index.html'),
        learn: resolve(_dirname, 'learn/index.html'),
        ecosystem: resolve(_dirname, 'ecosystem/index.html'),
        docs: resolve(_dirname, 'docs/index.html')
      },
    },
  },
});
