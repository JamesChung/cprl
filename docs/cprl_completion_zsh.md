## cprl completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(cprl completion zsh); compdef _cprl cprl

To load completions for every new session, execute once:

#### Linux:

	cprl completion zsh > "${fpath[1]}/_cprl"

#### macOS:

	cprl completion zsh > $(brew --prefix)/share/zsh/site-functions/_cprl

You will need to start a new shell for this setup to take effect.


```
cprl completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --aws-profile string   overrides [aws-profile] value in cprl.yaml
      --gov-cloud            set context as gov-cloud
      --profile string       references a profile in cprl.yaml (default "default")
```

### SEE ALSO

* [cprl completion](cprl_completion.md)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 20-Mar-2023