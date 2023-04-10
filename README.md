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

> This will build and install the binary into your `~/.local/bin`. Make sure you have `~/.local/bin` in your `PATH`.

```sh
$ make local
```

## Handy Examples

> `cprl credentials assume` will prompt you for a role ARN and session name and will configure an AWS profile based on the name you provide it (or you could provide the same information via flags). You can combine `cprl credentials assume` with `cprl credentials output` and make that AWS profile your current active session.

```sh
$ cprl --aws-profile=main credentials assume --role-arn=arn:aws:iam::010203040506:role/dev --session-name=cprl --output-profile=dev
$ source <(cprl credentials output --aws-profile=dev)
$ aws sts get-caller-identity
{
    "UserId": "TAG0YY70NST6IUO5KA5XB:cprl",
    "Account": "010203040506",
    "Arn": "arn:aws:sts::010203040506:assumed-role/dev/cprl"
}
```

> `cprl console open --aws-profile=dev` will open your default web browser to the AWS console based on the specified AWS profile.

## Commands

- [`codecommit`](./docs/cprl_codecommit.md)
  - [`branch`](./docs/cprl_codecommit_branch.md)
    - [`remove`](./docs/cprl_codecommit_branch_remove.md)
  - [`pr`](./docs/cprl_codecommit_pr.md)
    - [`approve`](./docs/cprl_codecommit_pr_approve.md)
    - [`closes`](./docs/cprl_codecommit_pr_close.md)
    - [`create`](./docs/cprl_codecommit_pr_create.md)
    - [`diff`](./docs/cprl_codecommit_pr_diff.md)
    - [`list`](./docs/cprl_codecommit_pr_list.md)
- [`console`](./docs/cprl_console.md)
  - [`open`](./docs/cprl_console_open.md)
- [`credentials`](./docs/cprl_credentials.md)
  - [`assume`](./docs/cprl_credentials_assume.md)
  - [`clear`](./docs/cprl_credentials_clear.md)
  - [`output`](./docs/cprl_credentials_output.md)
  - [`list`](./docs/cprl_credentials_list.md)

**The documentation for each command and sub-command can be found in [./docs](./docs/cprl.md).**

### Help with auto-completions

* [bash completion](./docs/cprl_completion_bash.md)
* [fish completion](./docs/cprl_completion_fish.md)
* [zsh completion](./docs/cprl_completion_zsh.md)
* [powershell completion](./docs/cprl_completion_powershell.md)

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
