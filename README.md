# mcli

Personal CLI for managing your own PCs/servers: Wake-on-LAN, SSH shortcuts, PC health checks, Docker cleanup, system updates, Wi-Fi hotspot control, network status, and quick file sharing — all in one Go binary with an interactive TUI and scriptable subcommands.

Runs on Linux and Android/Termux.

## Install

**Via Go:**

```bash
go install github.com/skygrime35/mcli@latest
```

**Via install script:**

```bash
curl -fsSL https://raw.githubusercontent.com/skygrime35/mcli/main/scripts/install.sh | sh
```

Release binaries are published for `linux/amd64` and `linux/arm64` (covering desktop Linux and Termux/Android).

## Configuration

On first run, `mcli` creates `~/.config/mcli/config.yaml`. See [`config.example.yaml`](./config.example.yaml) for the full format, or add a server directly:

```bash
mcli server add --name home --host my-server.example.com --mac AA:BB:CC:DD:EE:FF --ssh-user myuser
```

## Usage

Run `mcli` with no arguments for the interactive menu, or use subcommands directly (handy for scripts/cron):

```bash
mcli server list
mcli server wol home
mcli server ssh home
```

Run `mcli --help` for the full command list.

## Status

Feature-complete baseline (v0.1.0): this is a from-scratch Go rewrite of an earlier Python version. All planned features are implemented — interactive TUI (bubbletea), Server Manager, PC Health (CLI + TUI, summary and watch modes), Docker Manager, System Update, Hotspot Manager, Network Status, and File Sharing. Every main-menu item is a working feature with no remaining placeholders.

As a personal tool, it prioritizes your own PC/server workflows over broad platform coverage or edge-case stability guarantees. It runs on Linux (desktop) and Android/Termux.
