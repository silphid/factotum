# Commands

- factotum use [context or "none"]
  - validates context
  - saves context in `~/.factotum/state.yaml`
  - omitting context prompts for context from list of available contexts + "none"
- factotum start [context or "none"]
  - omitting context uses current context or prompts for context
  - reads `~/.factotum/state.yaml`
  - if no version set yet, does a `factotum upgrade`
  - merges env vars and volume mounts from `shared.yaml` and `user.yaml`
  - keep only volume mounts that exist locally
  - starts docker container
- factotum stop [context or "all"]
  - omitting context prompts for context
- factotum list
  - contexts
    - lists available contexts
  - containers
    - lists active containers
- factotum pull
  - pulls latest factotum config git repo
- factotum remove
  - stops all running containers
  - deletes factotum docker images and containers
- factotum upgrade
  - discovers latest version of factotum docker image
  - stores latest version number in `~/.factotum/state.yaml`
  - docker pull latest version

## Global flags

- -h --help
- -v --verbose
  - Displays detailed output messages
- --dry-run
  - Only displays shell commands that it would execute normally
  - Automatically turns-on verbose mode

# Config files structure

- ~/.factotum/
  - settings.yaml
  - user.yaml (user-specific overrides)
- Factotum Git Repo Clone Dir
  - config
    - shared.yaml (company-wide base values)
    - user.yaml (default user configs, copied to home during install)

# ~/.factotum/state.yaml file format

```yaml
clone: /Users/mathieu/dev/factotum
version: 1.2.3
context: cluster1
```

# Config files format

```yaml
container:
  registry: dockerhub # supported values are `gcr`, `ecr` and `dockerhub`
  image: silphid/factotum

contexts:
  - name: cluster1
    env:
      KUBE_CONTEXT: cluster1
      # REGION: us-east-2
  - name: cluster2
    env:
      KUBE_CONTEXT: cluster2

env:
  CLOUD: aws # supported clouds: aws, gcp
  REGION: us-east-1

volumes:
  $HOME/.ssh: /root/.ssh
  $HOME/.gitconfig: /root/.gitconfig
  $HOME/.aws: /root/.aws
  $HOME/.config/gh: /root/.config/gh
  $HOME/.cfconfig: /root/.cfconfig
```

# Installation

## Factotum git repo structure

- cli
  - go cli source code
- docker
  - Dockerfile + image source files
- config
  - user.yaml (default user-specific config copied to ~/.factotum/user.yaml, if not already existing)
  - shared.yaml (shared across organisation/teams)
- install.sh

## Process

- User clones factotum git repo to folder where it will permanently reside
- From repo root, user runs `./install.sh`, which does:
  - Copies `./config/user.yaml` to `~/.factotum/user.yaml`
  - Creates `~/.factotum/
  - Downloads tar.gz file for latest build of cli for current OS and architecture
  - Decompresses and copies to /usr/local/bin
  - Runs `factotum upgrade`
