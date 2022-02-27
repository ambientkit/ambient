# Ambient core <!-- omit in toc -->

[![GitHub Actions status](https://github.com/ambientkit/ambient/actions/workflows/unit-tests.yml/badge.svg)](https://github.com/ambientkit/ambient/actions)

Thanks for visiting! All docs are available [here](https://ambientkit.github.io/).

This repository contains the code for the Ambient core which includes:

- [entry points](https://github.com/ambientkit/ambient/blob/main/pkg/ambientapp/app.go) like `NewApp()`, `NewAppLogger()` and `ListenAndServe()`
- interfaces used by the plugins: [logger](https://github.com/ambientkit/ambient/blob/main/ambient_logger.go), [router](https://github.com/ambientkit/ambient/blob/main/ambient_router.go), [session manager](https://github.com/ambientkit/ambient/blob/main/ambient_session.go), etc.
- [permission system](https://github.com/ambientkit/ambient/tree/main/internal/secureconfig) and [permission list](https://github.com/ambientkit/ambient/blob/main/grant.go)
- [dev console](https://github.com/ambientkit/ambient/blob/main/internal/devconsole/devconsole.go) for the [AMB CLI](https://github.com/ambientkit/amb)
- [toolkit](https://github.com/ambientkit/ambient/blob/main/toolkit.go) available to all plugins