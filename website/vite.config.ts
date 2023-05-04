import { defineConfig } from 'vite';
import { createHtmlPlugin } from 'vite-plugin-html'

export default defineConfig({
  plugins: [createHtmlPlugin({
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
  },
});
