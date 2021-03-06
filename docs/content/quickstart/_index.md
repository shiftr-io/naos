---
title: Quickstart
---

# Quickstart

The following steps illustrate how you can get started with NAOS. We will guide you through installation, project setup and flashing.

## Installation

First, you need to install the latest version of the `naos` command line utility (CLI):

```
curl -L https://naos.256dpi.com/install.sh | bash
```

*You can also download the binary manually from <https://github.com/256dpi/naos/releases> if you don't want to run the above shell script.*

After the installation you can verify that `naos` is available:

```
naos help
```

## Project Setup

Create a new project in an empty directory somewhere on your computer:

```
naos create
```

*The CLI will create a `naos.json` configuration and an empty `src/main.c`  file.*

Download and install all dependencies (this may take a few minutes):

```
naos install
```

*You can run `naos install` anytime to update to the dependencies.*

Run the firmware on the connected ESP32:

```
naos run
```
