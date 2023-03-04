# cprl

## Config File

`cprl` will first search for a `cprl.yaml` file in the current working directory. If not found it will search in the user's home `.config/` directory as `.config/cprl/cprl.yaml`. If neither is found `cprl` will prompt you if you'd like it to create a `cprl.yaml` file for you with a template. You will still need to provide it your preferred values after creation.

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
