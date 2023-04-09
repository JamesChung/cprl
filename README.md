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

## Handy Examples

> `cprl credentials assume` will prompt you for a role ARN and session name and will configure an AWS profile based on the name you provide it.

> `cprl console open --aws-profile=dev` will open your default web browser to the AWS console based on the specified AWS profile.

## Commands

**The documentation for each command and sub-command can be found [here](./docs/cprl.md).**

## Config File

> `cprl` will first search for a `cprl.yaml` file in the current working directory. If not found it will search in the user's home `.config/` directory as `.config/cprl/cprl.yaml`. If neither is found `cprl` will prompt you if you'd like it to create a `cprl.yaml` file for you with a template. You will still need to provide it your preferred values after creation.

### Schema

```yaml
# cprl will always default to this profile if `--profile` is not set
default:
  # profile wide configs
  config:
    # the default aws profile used unless overrode via `--aws-profile` flag
    aws-profile: <profile name>
  # individual service level configurations
  services:
    console:
      gov-cloud: true | false
    codecommit:
      repositories:
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
