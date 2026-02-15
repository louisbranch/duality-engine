---
title: "Docs Site Organization"
parent: "Notes"
nav_order: 2
---

# Docs Site Organization

## Navigation rules

- Every page should set `title` and `parent` in front matter.
- Section indexes use `has_children: true` and a `nav_order` that keeps the sidebar stable.
- Generated docs should emit their own front matter so they stay grouped.

## Public build

The public site is limited to Audience, Running, Project, and Reference.
Use the public config to exclude internal sections and internal-only Project docs.

Build command:

`jekyll build --config _config.yml,_config.public.yml`

## Public exclusions

`docs/_config.public.yml` excludes these directories:

- `events`
- `guides`
- `specs`
- `product`
- `notes`

It also excludes internal Project pages:

- `project/identity.md`
- `project/oauth.md`
- `project/participant-invitation.md`
- `project/testing-scenarios.md`
- `project/scenario-dsl-dependencies.md`
- `project/scenario-missing-mechanics.md`
- `project/icon-catalog.md`
