/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        './index.html',
        './docs/index.html',
        './learn/index.html',
        './ecosystem/index.html',
        './src/**/*.{js,ts,jsx,tsx,css,md,mdx,html,json,scss}',
    ],
    darkMode: 'class',
    theme: {
        extend: {},
    },
    plugins: [
        require('@tailwindcss/typography'),
    ],
};