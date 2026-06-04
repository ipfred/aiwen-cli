#!/usr/bin/env node

const fs = require("fs");
const path = require("path");
const { execFileSync, execFile } = require("child_process");
const p = require("@clack/prompts");

const PKG = "aiwen-geoip-cli";
const SKILLS_REPO = "ipfred/aiwen-skills";
const isWindows = process.platform === "win32";

// ---------------------------------------------------------------------------
// i18n
// ---------------------------------------------------------------------------

const messages = {
  zh: {
    setup:          "正在设置 aw-cli...",
    step1:          "正在安装 %s...",
    step1Upgrade:   "正在升级 %s (v%s → v%s)...",
    step1Skip:      "已安装 (v%s)，跳过",
    step1Done:      "已全局安装",
    step1Upgraded:  "已升级到 v%s",
    step1Fail:      "全局安装失败。运行以下命令重试: npm install -g %s",
    step2:          "配置 API Key",
    step2Prompt:    "请输入你的 IPPlus360 / AIWEN API Key",
    step2Skip:      "跳过 API Key 配置",
    step2Done:      "API Key 已配置",
    step2Fail:      "API Key 配置失败。运行以下命令重试: aw-cli config set api_key YOUR_KEY",
    step3:          "安装 AI Skills",
    step3Spinner:   "正在安装 Skills...",
    step3Skip:      "已安装，跳过",
    step3Done:      "Skills 已安装",
    step3Fail:      "Skills 安装失败。运行以下命令重试: npx skills add %s -y -g",
    step4:          "验证安装",
    step4Done:      "安装验证完成",
    step4Fail:      "验证失败，请检查安装",
    done:           "安装完成！你可以试试：\n  aw-cli loc 8.8.8.8\n  aw-cli current\n  aw-cli --help",
    cancelled:      "安装已取消",
    nonTtyHint:     "要完成配置，请在终端中运行：\n  aw-cli config set api_key YOUR_KEY\n  aw-cli --help",
  },
  en: {
    setup:          "Setting up aw-cli...",
    step1:          "Installing %s globally...",
    step1Upgrade:   "Upgrading %s (v%s → v%s)...",
    step1Skip:      "Already installed (v%s). Skipped",
    step1Done:      "Installed globally",
    step1Upgraded:  "Upgraded to v%s",
    step1Fail:      "Failed to install globally. Run manually: npm install -g %s",
    step2:          "Configure API Key",
    step2Prompt:    "Enter your IPPlus360 / AIWEN API Key",
    step2Skip:      "Skipped API Key configuration",
    step2Done:      "API Key configured",
    step2Fail:      "Failed to configure API Key. Run: aw-cli config set api_key YOUR_KEY",
    step3:          "Install AI skills",
    step3Spinner:   "Installing skills...",
    step3Skip:      "Already installed. Skipped",
    step3Done:      "Skills installed",
    step3Fail:      "Failed to install skills. Run manually: npx skills add %s -y -g",
    step4:          "Verify installation",
    step4Done:      "Installation verified",
    step4Fail:      "Verification failed, please check your installation",
    done:           "You're all set! Try:\n  aw-cli loc 8.8.8.8\n  aw-cli current\n  aw-cli --help",
    cancelled:      "Installation cancelled",
    nonTtyHint:     "To complete setup, run:\n  aw-cli config set api_key YOUR_KEY\n  aw-cli --help",
  },
};

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function handleCancel(value, msg) {
  if (p.isCancel(value)) {
    p.cancel(msg.cancelled);
    process.exit(0);
  }
  return value;
}

function execCmd(cmd, args, opts) {
  if (isWindows) {
    return execFileSync("cmd.exe", ["/c", cmd, ...args], opts);
  }
  return execFileSync(cmd, args, opts);
}

function run(cmd, args, opts = {}) {
  execCmd(cmd, args, { stdio: "inherit", ...opts });
}

function runSilent(cmd, args, opts = {}) {
  return execCmd(cmd, args, {
    stdio: ["ignore", "pipe", "pipe"],
    ...opts,
  });
}

function runSilentAsync(cmd, args, opts = {}) {
  const actualCmd = isWindows ? "cmd.exe" : cmd;
  const actualArgs = isWindows ? ["/c", cmd, ...args] : args;
  return new Promise((resolve, reject) => {
    execFile(actualCmd, actualArgs, {
      stdio: ["ignore", "pipe", "pipe"],
      ...opts,
    }, (err, stdout) => {
      if (err) reject(err);
      else resolve(stdout);
    });
  });
}

function fmt(template, ...values) {
  let i = 0;
  return template.replace(/%s/g, () => values[i++] ?? "");
}

/** Resolve the path of globally installed aw-cli (skip npx temp copies). */
function whichAwCli() {
  try {
    const prefix = execFileSync("npm", ["prefix", "-g"], {
      stdio: ["ignore", "pipe", "pipe"],
    }).toString().trim();
    const bin = isWindows
      ? path.join(prefix, "aw-cli.cmd")
      : path.join(prefix, "bin", "aw-cli");
    if (fs.existsSync(bin)) return bin;
  } catch (_) {
    // fall through
  }
  // Fallback to which/where if npm prefix lookup fails.
  try {
    const cmd = isWindows ? "where" : "which";
    return execFileSync(cmd, ["aw-cli"], { stdio: ["ignore", "pipe", "pipe"] })
      .toString()
      .split("\n")[0]
      .trim();
  } catch (_) {
    return null;
  }
}

/** Get the latest version of aiwen-geoip-cli from the registry. Returns version or null. */
function getLatestVersion() {
  try {
    const out = runSilent("npm", ["view", PKG, "version"], { timeout: 15000 });
    const ver = out.toString().trim();
    return /^\d+\.\d+\.\d+/.test(ver) ? ver : null;
  } catch (_) {
    return null;
  }
}

/** Compare two semver strings. Returns true if a < b. */
function semverLessThan(a, b) {
  const pa = a.replace(/-.*$/, "").split(".").map(Number);
  const pb = b.replace(/-.*$/, "").split(".").map(Number);
  for (let i = 0; i < 3; i++) {
    if ((pa[i] || 0) < (pb[i] || 0)) return true;
    if ((pa[i] || 0) > (pb[i] || 0)) return false;
  }
  return false;
}

/** Check whether aiwen-geoip-cli is truly installed in npm global prefix. Returns version or null. */
function getGloballyInstalledVersion() {
  try {
    const out = runSilent("npm", ["list", "-g", PKG], { timeout: 15000 });
    const match = out.toString().match(/@(\d+\.\d+\.\d+[^\s]*)/);
    return match ? match[1] : "unknown";
  } catch (_) {
    return null;
  }
}

/** Check whether aw-cli config has an API key set. */
function hasApiKey(binPath) {
  try {
    const out = runSilent(binPath, ["config", "show"], { timeout: 10000 });
    const json = JSON.parse(out.toString());
    return !!(json.api_key || json.apiKey);
  } catch (_) {
    return false;
  }
}

/** Parse --lang from process.argv, returns "zh", "en", or null. */
function parseLangArg() {
  const args = process.argv.slice(2);
  for (let i = 0; i < args.length; i++) {
    if (args[i] === "--lang" && args[i + 1]) {
      const val = args[i + 1].toLowerCase();
      if (val === "zh" || val === "en") return val;
    }
    if (args[i].startsWith("--lang=")) {
      const val = args[i].split("=")[1].toLowerCase();
      if (val === "zh" || val === "en") return val;
    }
  }
  return null;
}

// ---------------------------------------------------------------------------
// Steps
// ---------------------------------------------------------------------------

async function stepSelectLang() {
  const fromArg = parseLangArg();
  if (fromArg) return fromArg;

  const lang = await p.select({
    message: "请选择语言 / Select language",
    options: [
      { value: "zh", label: "中文" },
      { value: "en", label: "English" },
    ],
  });
  return handleCancel(lang, messages.zh);
}

async function stepInstallGlobally(msg) {
  const installedVer = getGloballyInstalledVersion();
  const latestVer = getLatestVersion();
  const needsUpgrade = installedVer && latestVer && semverLessThan(installedVer, latestVer);

  if (installedVer && !needsUpgrade) {
    p.log.info(fmt(msg.step1Skip, installedVer));
    return false;
  }

  const s = p.spinner();
  if (needsUpgrade) {
    s.start(fmt(msg.step1Upgrade, PKG, installedVer, latestVer));
  } else {
    s.start(fmt(msg.step1, PKG));
  }
  try {
    await runSilentAsync("npm", ["install", "-g", PKG], { timeout: 120000 });
    s.stop(needsUpgrade ? fmt(msg.step1Upgraded, latestVer) : msg.step1Done);
    return needsUpgrade;
  } catch (_) {
    s.stop(fmt(msg.step1Fail, PKG));
    process.exit(1);
  }
}

async function stepConfigApiKey(msg) {
  const awCli = whichAwCli();
  if (!awCli) {
    p.log.warn("aw-cli not found, skipping API key configuration");
    return;
  }

  if (hasApiKey(awCli)) {
    p.log.info(msg.step2Skip);
    return;
  }

  const apiKey = await p.password({
    message: msg.step2Prompt,
  });
  if (p.isCancel(apiKey)) {
    p.cancel(msg.cancelled);
    process.exit(0);
  }

  if (!apiKey) {
    p.log.info(msg.step2Skip);
    return;
  }

  try {
    run(awCli, ["config", "set", "api_key", apiKey]);
    p.log.success(msg.step2Done);
  } catch (_) {
    p.log.error(msg.step2Fail);
  }
}

async function skillsAlreadyInstalled() {
  try {
    const out = await runSilentAsync("npx", ["-y", "skills", "ls", "-g"], {
      timeout: 120000,
    });
    return /^aw-cli-/m.test(out.toString());
  } catch (_) {
    return false;
  }
}

async function stepInstallSkills(msg) {
  const s = p.spinner();
  s.start(msg.step3Spinner);
  try {
    if (await skillsAlreadyInstalled()) {
      s.stop(msg.step3Skip);
      return;
    }
    try {
      await runSilentAsync("npx", ["-y", "skills", "add", SKILLS_REPO, "-y", "-g"], {
        timeout: 120000,
      });
    } catch (primaryErr) {
      s.stop(fmt(msg.step3Fail, SKILLS_REPO));
      process.exit(1);
    }
    s.stop(msg.step3Done);
  } catch (_) {
    s.stop(fmt(msg.step3Fail, SKILLS_REPO));
    process.exit(1);
  }
}

async function stepVerify(msg) {
  const awCli = whichAwCli();
  if (!awCli) {
    p.log.warn(msg.step4Fail);
    return;
  }

  try {
    const out = runSilent(awCli, ["--version"], { timeout: 10000 });
    p.log.success(`${msg.step4Done}: ${out.toString().trim()}`);
  } catch (_) {
    p.log.warn(msg.step4Fail);
  }
}

// ---------------------------------------------------------------------------
// Main
// ---------------------------------------------------------------------------

async function main() {
  const isInteractive = !!process.stdin.isTTY;
  const lang = isInteractive ? await stepSelectLang() : (parseLangArg() || "en");
  const msg = messages[lang];

  if (isInteractive) {
    p.intro(msg.setup);
    await stepInstallGlobally(msg);
    await stepConfigApiKey(msg);
    await stepInstallSkills(msg);
    await stepVerify(msg);
    p.outro(msg.done);
  } else {
    console.log(msg.setup);
    await stepInstallGlobally(msg);
    await stepInstallSkills(msg);
    console.log(msg.nonTtyHint);
  }
}

main().catch((err) => {
  p.cancel("Unexpected error: " + (err.message || err));
  process.exit(1);
});
