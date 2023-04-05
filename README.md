# [WARNING] In Early Development

# cprl

## How to Install

### Homebrew

```sh
$ brew tap JamesChung/tap
$ brew install JamesChung/tap/cprl
```

### Go

```sh
$ go install github.com/JamesChung/cprl@latest
```

### Make

> This will build and install the binary to your `~/.local/bin`. Make sure you have that in your `PATH`.

```sh
$ make local
```

## Commands

**[Documentation link](./docs/cprl.md)**

## Config File

> `cprl` will first search for a `cprl.yaml` file in the current working directory. If not found it will search in the user's home `.config/` directory as `.config/cprl/cprl.yaml`. If neither is found `cprl` will prompt you if you'd like it to create a `cprl.yaml` file for you with a template. You will still need to provide it your preferred values after creation.

### Schema

```yaml
default:                            # cprl will always default to this profile
  config:                           # profile wide configs
    aws-profile: <profile name>     # this aws profile will be used by default for commands
  services:                         # individual service level configurations
    console:
      gov-cloud: true | false
    codecommit:                     # name of a supported service
      repositories:                 # service specific configurations
        - <repo name>
        - <repo name>
        - <repo name>
<custom profile>:
  config:
    aws-profile: <profile name>
  services:
    codecommit:
      repositories:
        - <repo name>
        - <repo name>
        - <repo name>
```

### Example

```yaml
default:
  config:
    aws-profile: default
  services:
    codecommit:
      repositories:
        - example-repo
        - other-example-repo
secondary:
  config:
    aws-profile: dev
  services:
    codecommit:
      repositories:
        - dev-example-repo
```
