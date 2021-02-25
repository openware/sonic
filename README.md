# Sonic Fullstack micro-framework

Sonic is a project templates for creating server-side rendered applications. Powered by [gin](https://github.com/gin-gonic/gin)

## Roadmap

  - Integrate with Baseapp
  - CMS for dynamic pages

### Repo structure

1. `scripts` - scripts for generating & updating your application.
2. `skel` - a skeleton for your app.
3. `skel/config` - application config files.
4. `skel/handlers` - REST handlers for CMS.
5. `skel/models` - models for database entities.

## How to generate an app

```bash
curl -ssL https://raw.githubusercontent.com/openware/sonic/master/scripts/install.sh | zsh
svm create github.com/*username*/*project_name*
```

## Setup

Setup database:

```
go run . db create
go run . db migrate
```

Run server:

```
go run . serve
```
