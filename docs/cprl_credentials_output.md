## cprl credentials output

outputs AWS credentials

```
cprl credentials output [flags]
```

### Examples

```
  Basic output:
  $ cprl credentials output
  export AWS_ACCESS_KEY_ID=<access key id value>
  export AWS_SECRET_ACCESS_KEY=<secret access key value>
  export AWS_SESSION_TOKEN=<session token value>
  
  JSON output:
  $ cprl credentials output --json
  {"AccessKeyID":"<access key id value>","SecretAccessKey":...}
  
  Source credentials as your current session:
  $ source <(cprl credentials output --aws-profile=dev)
  $ aws sts get-caller-identity
  {
  "UserId": "TAG0YY70NST6IUO5KA5XB:cprl",
  "Account": "010203040506",
  "Arn": "arn:aws:sts::010203040506:assumed-role/dev/cprl"
  }
```

### Options

```
  -h, --help   help for output
      --json   output in json format
```

### Options inherited from parent commands

```
      --aws-profile string   overrides [aws-profile] value in cprl.yaml
      --profile string       references a profile in cprl.yaml (default "default")
```

### SEE ALSO

* [cprl credentials](cprl_credentials.md)	 - AWS credentials

###### Auto generated by spf13/cobra on 9-Apr-2023