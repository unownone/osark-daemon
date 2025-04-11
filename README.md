# OSARK-Daemon: Comprehensive Tracking daemon for OSARK

This daemon is responsible for tracking the user's activity on the system.

## Features

- Tracking
  - [x] Apps & System Activity Info
  - [ ] App Open/Close/ Focus Events

- Reporting
  - [x] Pushing reports to the server
  - [x] Batched reporting for efficient network usage
  - [ ] kafka/ queue based pushing to avoid latency and scale-up

- Configuration
  - [x] Daemon process
  - [ ] Auto Startup
  - [ ] Auto installation of osquery if not present

## Installation

```bash
go install github.com/unownone/osark-daemon@latest
```

## Usage

```bash
make build;
./osark-daemon
```
