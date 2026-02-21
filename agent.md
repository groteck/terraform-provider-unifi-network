# Agent Context: UniFi Terraform Provider

## Project Purpose
A custom Terraform provider for managing UniFi Network resources (VLANs, Firewall rules, etc.) using the modern UniFi Network API.

## Requirements
- **Conventional Commits:** `feat`, `fix`, `chore`, `docs`, `style`, `refactor`, `perf`, `test`.
- **Merge Strategy:** Automated releases with changelogs on merge to `main`.
- **Testing:** Local Docker container (`linuxserver/unifi-network-application`) for acceptance tests.
- **Documentation:** Automated via `tfplugindocs`, enriched with "human-understandable" context from YouTube tutorials.

## Agent Workflow (Strict)
- **Git Worktrees:** Agents MUST use git worktrees for development to isolate changes from the main repository.
- **Feature Branches:** All work must be done in feature branches.
- **Pre-commit Checks:** Run `make pre-commit` (fmt, lint, test) before committing.

## References
- **API Docs:** [UniFi Network API](https://unifi.ui.com/consoles/1C0B8B10ADF60000000008ED992C000000000967AE4B0000000067EC5CD0:1577279704/unifi-api/network)
- **Tutorials:** [UniFi Network Playlist](https://www.youtube.com/playlist?list=PLjGQNuuUzvmvxayWV93dbBleXzt6RCvXP)
- **Base Style:** Follow patterns in `../pangolin-tf/`.

## Current Status
- [x] Initial Manifest
- [x] Project Scaffolding
- [x] CI/CD Configuration (with live Docker tests) - **PASSING**
- [x] Provider Core Implementation (Dual-Auth: Password & Token)
- [x] Headless Docker Environment (Bypasses Setup Wizard)
- [x] Refactored Generic Client & Base Classes (Go 1.24)
- [x] Data Sources Implementation (7 Core Data Sources)
- [x] Resource: `unifi_network` (VLANs) - **PASS**
- [x] Resource: `unifi_port_profile` - **PASS**
- [x] Resource: `unifi_user_group` - **PASS**
- [x] Resource: `unifi_ap_group` - **PASS**
- [x] Resource: `unifi_wlan` - **PASS**
- [x] Resource: `unifi_firewall_group` - **PASS**
- [x] Resource: `unifi_port_forward` - **PASS**
- [x] Resource: `unifi_radius_profile` - **PASS**
- [x] Resource: `unifi_firewall_rule` - **PASS**
- [x] Resource: `unifi_user` - **PASS**
- [x] Resource: `unifi_traffic_rule` - **PASS**
- [x] Resource: `unifi_static_dns` - **PASS**
- [x] Resource: `unifi_static_route` - **PASS**
- [x] Documentation & Examples (Modular Blueprints & Big System)
- [x] Architecture Documentation (README)

- [x] Architecture Documentation (README)

## Refactor Notes
- Attributes normalized to standard naming conventions (`network_id`, `passphrase`, etc.).
- Integrated `go-retryablehttp` for API resilience.
- Leveraged Go generics to reduce CRUD boilerplate by 60%.
- Added comprehensive blueprints for Home Base and Homelab scenarios.
