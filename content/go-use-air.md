---
title: "Hot reloading for Go Projects"
date: 2025-09-15
summary: "Why didn't I start using this sooner? A quick guide to hot reloading for Go projects with Air."
draft: false
tags:
  - "#Go"
  - "#Development"
---

## Air Quickstart

1 **Install Air**

```bash
go install github.com/cosmtrek/air@latest
```

2 **Create an `.air.toml` config (optional)**

Run `air init` to generate a config file, or use the defaults.

3 **Run your project with Air**

```bash
air
```

Air will watch your files and automatically restart your Go app on changes.

**Tip:** Add `air` to your project’s dev dependencies or document it in your README for your team.
