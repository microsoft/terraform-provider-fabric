# Local Terraform Provider Development configuration

## DevContainer usage

Add to your `devcontainer.json`:

```json
"features": {
  "./features/tfprovider-local-dev": {
    "providerName": "<your-owner>/<your-provider-name>",
    "workspace": "${containerWorkspaceFolder}"
  }
}
```

## WSL installation

Under root of your repo, in the console set your provider name and run:

```shell
export PROVIDERNAME="<org>/<provider>"
.devcontainer/features/tfprovider-local-dev/install.sh
```
