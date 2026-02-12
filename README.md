# Padel Services

This repository contains a set of services geared
towards Padel clubs.
Currently the following services are offered/planned:

- [x] printing scoreboard for Rodeo tournament
- [x] match-making for Rodeo tournament
- [ ] tournament live updates
- [ ] support for more kinds of tournaments (e.g. direct-elimination)

Services are provided via a telegram mini-app.

## Architectural details

Go web-server that serves a single page application and exposes API for
further interaction.

The frontend uses React and React Router

- [Gin](https://gin-gonic.com/) for the web framework.
- The frontend uses React and React Router, and [Telegram UI](https://github.com/telegram-mini-apps-dev/TelegramUI), which provides pre-baked
components, stylized to provide a native-feel interaction with the mini-app.
A third party javascript SDK [tma.js](https://docs.telegram-mini-apps.com/),
simplifies interactions with Telegram.
- Persistency layer implemented via PostgreSQL ([schema](resources/sql/create_tables.sql)).
