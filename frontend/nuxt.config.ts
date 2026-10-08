export default defineNuxtConfig({
  compatibilityDate: "2026-10-06",
  ssr: false,

  modules: ["@nuxt/ui", "@nuxt/icon"],

  css: ["~/assets/css/main.css", "leaflet/dist/leaflet.css"],

  // Отключаем загрузку шрифтов
  fonts: {
    providers: {
      google: false,
      fontshare: false,
      bunny: false,
      adobe: false,
      googleicons: false,
    },
  },

  colorMode: {
    preference: "dark",
    fallback: "dark",
  },

  runtimeConfig: {
    public: {
      apiBase: "http://picmap.ru",
    },
  },

  app: {
    head: {
      htmlAttrs: { lang: "ru" },
      title: "PicMap",
      meta: [
        {
          name: "description",
          content: "Твои путешествия. Твои воспоминания.",
        },
        { name: "viewport", content: "width=device-width, initial-scale=1.0" },
      ],
    },
  },

  devtools: { enabled: true },
});
