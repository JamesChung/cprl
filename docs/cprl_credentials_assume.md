## cprl credentials assume

assume AWS role

```
cprl credentials assume [flags]
```

### Examples

```
  Assume role:
  $ cprl credentials assume
  ...
  
  Assume role bypassing input prompts via flags:
  $ cprl --aws-profile=main credentials assume \
  --role-arn=arn:aws:iam::010203040506:role/dev \
  --session-name=cprl --output-profile=dev
```

### Options

```
  -h, --help                    help for assume
      --output-profile string   new profile name of the assuming role
      --role-arn string         role ARN of the assuming role
      --session-name string     name of the session
```

### Options inherited from parent commands

```
      --aws-profile string   overrides [aws-profile] value in cprl.yaml
      --profile string       references a profile in cprl.yaml (default "default")
```

### SEE ALSO

* [cprl credentials](cprl_credentials.md)	 - AWS credentials

###### Auto generated by spf13/cobra on 27-Feb-2024
