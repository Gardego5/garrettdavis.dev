/** @type {import('tailwindcss').Config} */
module.exports = {
  content: {
    relative: true,
    files: ["**/*.{go,md}"],
  },
  theme: {
    fontFamily: {
      mono: ["Iosevka GarrettDavisDev Web", "ui-monospace"],
    },
    extend: {},
  },
  plugins: [],
};
