# nMeshed Repo Migration Guide

## Goal
Split codebase into **public** (SDKs, examples, docs) and **private** (sync server IP) repos.

---

## Phase 1: GitHub Setup ✅

### 1.1 Create GitHub Organization
- [X] Go to [github.com/organizations/plan](https://github.com/organizations/plan)
- [X] Click "Create free organization"
- [X] **Name**: `nmeshed`
- [X] Set as PUBLIC organization
- [X] Skip inviting members

### 1.2 Create Public Repo
- [X] In `nmeshed` org, click "New repository"
- [X] **Name**: `nmeshed`
- [X] **Visibility**: Public
- [X] **NO** README (we'll push our own)
- [X] Click "Create repository"

### 1.3 Create Private Repo
- [X] In `nmeshed` org, click "New repository"
- [X] **Name**: `sync-server`
- [X] **Visibility**: Private
- [X] **NO** README
- [X] Click "Create repository"

---

## Phase 2: Prepare Local Folders ✅

```bash
# Create new project directory
mkdir -p ~/projects/nmeshed
cd ~/projects/nmeshed

# Clone empty repos
git clone git@github.com:nmeshed/nmeshed.git
git clone git@github.com:nmeshed/sync-server.git
```

---

## Phase 3: Copy PUBLIC Code (In Progress)

Copy these folders from `etherply` → `nmeshed/nmeshed`:

```bash
# From your old repo
cd ~/projects/etherply

# Create directories first
mkdir -p ~/projects/nmeshed/nmeshed/packages
mkdir -p ~/projects/nmeshed/nmeshed/pkg
mkdir -p ~/projects/nmeshed/nmeshed/apps

# Copy SDKs
cp -r packages/sdk-js ~/projects/nmeshed/nmeshed/packages/
cp -r packages/sdk-python ~/projects/nmeshed/nmeshed/packages/
cp -r pkg/go-sdk ~/projects/nmeshed/nmeshed/pkg/

# Copy examples, apps, docs
cp -r examples ~/projects/nmeshed/nmeshed/
cp -r apps/web ~/projects/nmeshed/nmeshed/apps/
cp -r apps/docs ~/projects/nmeshed/nmeshed/apps/
cp -r docs ~/projects/nmeshed/nmeshed/

# Copy root files
cp README.md ~/projects/nmeshed/nmeshed/
cp LICENSE ~/projects/nmeshed/nmeshed/
cp .gitignore ~/projects/nmeshed/nmeshed/
```

### Push public repo
```bash
cd ~/projects/nmeshed/nmeshed
git add .
git commit -m "Initial commit: SDKs, examples, docs"
git push origin main
```

---

## Phase 4: Copy PRIVATE Code

Copy the sync server (your IP):

```bash
cd ~/projects/etherply

# Copy sync server only (Go SDK is now public)
cp -r etherply-sync-server/* ~/projects/nmeshed/sync-server/

# Copy internal docs
cp -r internal-docs ~/projects/nmeshed/sync-server/
```

### Push private repo
```bash
cd ~/projects/nmeshed/sync-server
git add .
git commit -m "Initial commit: nMeshed Sync Server"
git push origin main
```

---

## Phase 5: Update Package References

### In public repo (`nmeshed/nmeshed`)
Update `packages/sdk-js/package.json`:
```json
{
  "repository": {
    "url": "https://github.com/nmeshed/nmeshed.git",
    "directory": "packages/sdk-js"
  },
  "homepage": "https://nmeshed.com"
}
```

### Commit and push
```bash
cd ~/projects/nmeshed/nmeshed
git add .
git commit -m "Update repo URLs"
git push origin main
```

---

## Phase 6: Archive Old Repo

- [ ] Go to `github.com/bneb/etherply`
- [ ] Settings → General → scroll to "Danger Zone"
- [ ] Click "Archive this repository"
- [ ] Type repo name to confirm

> **Don't delete yet!** Archive keeps it readable but prevents changes.

---

## Phase 7: Update Local Dev

Add both repos to your IDE:
```bash
# Open the new workspace
code ~/projects/nmeshed
```

Your new folder structure:
```
~/projects/nmeshed/
├── nmeshed/           # PUBLIC - SDKs, examples, docs
│   ├── packages/
│   ├── pkg/go-sdk/
│   ├── examples/
│   ├── apps/
│   └── docs/
└── sync-server/       # PRIVATE - Core engine
    ├── internal/
    ├── cmd/
    ├── internal-docs/
    └── go.mod
```

---

## Quick Reference

| What | Repo | Visibility |
|:---|:---|:---|
| JS SDK | `nmeshed/nmeshed` | Public |
| Python SDK | `nmeshed/nmeshed` | Public |
| Go SDK | `nmeshed/nmeshed` | Public |
| Examples | `nmeshed/nmeshed` | Public |
| Web App | `nmeshed/nmeshed` | Public |
| Docs | `nmeshed/nmeshed` | Public |
| Sync Server | `nmeshed/sync-server` | **Private** |
| Internal Docs | `nmeshed/sync-server` | **Private** |

---

## After Migration Checklist

- [ ] Verify `npm install nmeshed` still works
- [ ] Verify docs build (`npm run docs`)
- [ ] Test examples work
- [ ] Update Vercel to point to new repo (`nmeshed/nmeshed`)
- [ ] Update any CI/CD pipelines
