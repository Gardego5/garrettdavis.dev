const plugin = require("tailwindcss/plugin");

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: {
    relative: true,
    files: ["src/**/*.rs"],
  },
  theme: {
    fontFamily: {
      mono: ["Iosevka GarrettDavisDev Web", "ui-monospace"],
    },
    extend: {},
  },
  plugins: [],
};
