import { defineConfig } from "vite";

// https://vite.dev/config/
export default defineConfig({
  server: {
    allowedHosts: [`${process.env.REPO_NAME}.orb.local`],
  },
});
