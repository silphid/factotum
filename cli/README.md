# Commands

- factotum use [context or "base" or "none"]
  - validates context
  - saves context in `~/.factotum/state.yaml`
  - omitting context prompts for context from list of available contexts + "base" + "none"
- factotum start [context or "base"]
  - omitting context uses current context or prompts for context
  - reads `~/.factotum/state.yaml`
  - if no version set yet, does a `factotum upgrade`
  - merges env vars and volume mounts from `shared.yaml` and `user.yaml`
  - keep only volume mounts that exist locally
  - starts docker container
- factotum stop [context or "base" or "none" or "all"]
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
  - state.yaml
  - user.yaml: user-specific config overrides
- Factotum Git Repo Clone Dir
  - config
    - shared.yaml: base configs shared by all users
    - user.yaml: default user config file copied to home during install

# ~/.factotum/state.yaml file format

```yaml
version: 2021.04
cloneDir: /Users/mathieu/dev/factotum
imageVersion: 1.2.3
currentContext: cluster1
```

# Config files format

```yaml
base:
  registry: dockerhub # supported values are `gcr`, `ecr` and `dockerhub`
  image: silphid/factotum

  env:
    CLOUD: aws # supported clouds: aws, gcp
    REGION: us-east-1

  volumes:
    $HOME/.ssh: /root/.ssh
    $HOME/.gitconfig: /root/.gitconfig
    $HOME/.aws: /root/.aws
    $HOME/.config/gh: /root/.config/gh
    $HOME/.cfconfig: /root/.cfconfig

contexts:
  cluster1:
    env:
      KUBE_CONTEXT: cluster1
      # REGION: us-east-2
  cluster2:
    env:
      KUBE_CONTEXT: cluster2
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
  - Copies git repo `/config/user.yaml` to `~/.factotum/user.yaml`
  - Creates `~/.factotum/
  - Downloads tar.gz file for latest build of cli for current OS and architecture
  - Decompresses and copies to /usr/local/bin
  - Runs `factotum upgrade`

# Todo

- Add support for multiple images (other than factotum)
  - Sub-folders under `~/.factotum` and git `/config` that can specify extra
    `user.yaml` and `shared.yaml` files.
  - The sub-folder name can be used as a prefix to the context name (ie:
    `image1/context1`) or we add an extra `factotum image {image1}` command to
    select it?
