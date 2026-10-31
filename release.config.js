import fs from "node:fs/promises";
import path from "node:path";

export default {
  plugins: [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    {
      async prepare(_config, { cwd, nextRelease }) {
        const manifestPath = path.join(cwd, "herdr-plugin.toml");
        const manifest = await fs.readFile(manifestPath, "utf8");
        const versionLine = /^version = "[^"]*"$/gm;
        if ([...manifest.matchAll(versionLine)].length !== 1) {
          throw new Error("Expected exactly one version field in herdr-plugin.toml");
        }
        await fs.writeFile(
          manifestPath,
          manifest.replace(versionLine, `version = "${nextRelease.version}"`),
        );
      },
    },
    [
      "@semantic-release/git",
      {
        assets: ["herdr-plugin.toml"],
        message: "chore: bump herdr-plugin.toml to ${nextRelease.version}",
      },
    ],
    "@semantic-release/github",
  ],
};
