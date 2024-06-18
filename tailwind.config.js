const plugin = require("tailwindcss/plugin");

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.rs"],
  theme: {
    fontFamily: {
      mono: ["Iosevka Garrett Davis Dev Web", "ui-monospace"],
    },
    extend: {},
  },
  plugins: [
    plugin(function markdown({ addComponents, theme }) {
      /** @param {string} fontSize @param {string} lineHeight */
      const heading = (fontSize, lineHeight) => ({
        marginTop: theme("margin.4"),
        marginBottom: theme("margin.4"),
        fontSize: theme(`fontSize.${fontSize}`),
        lineHeight: theme(`lineHeight.${lineHeight}`),
        fontWeight: theme("fontWeight.semibold"),
        alignItems: "baseline",
        gap: "0.15em",
        position: "relative",
        display: "flex",
        marginLeft: "-0.5em",
        "&>a": {
          display: "inline-block",
          color: theme("colors.slate.400"),
          "&:hover": { color: theme("colors.blue.600") },
          scale: "0.8",
        },
      });

      addComponents({
        // @ts-ignore
        ".markdown": {
          "& p": {
            marginTop: theme("margin.3"),
            marginBottom: theme("margin.3"),
          },
          "& li": {
            listStyleType: "disc",
            listStylePosition: "outside",
            marginLeft: theme("margin.6"),
          },
          "& h1": heading("4xl", "10"),
          "& h2": heading("3xl", "9"),
          "& h3": heading("2xl", "8"),
          "& h4": heading("xl", "7"),
          "& h5": heading("lg", "7"),
          "& hr": {
            marginInline: theme("margin.6"),
            borderTop: `1px solid ${theme("borderColor.slate.500")}`,
          },

          lineHeight: theme("lineHeight.7"),
        },
      });
    }),
  ],
};
