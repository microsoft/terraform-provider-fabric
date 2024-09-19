# (optional, useful in the multi-tenant scenarios) Disable the new login experience
# See https://learn.microsoft.com/cli/azure/authenticate-azure-cli-interactively#sign-in-with-a-different-tenant for more details.
az config set core.login_experience_v2=off

# (optional, Windows only, useful in the multi-tenant scenarios) Disable WAM on Windows
# See https://learn.microsoft.com/cli/azure/authenticate-azure-cli-interactively#sign-in-with-web-account-manager-wam-on-windows for more details.
az config set core.enable_broker_on_windows=false

# Login to Azure with Entra ID credentials
# See https://learn.microsoft.com/cli/azure/authenticate-azure-cli for more details.
az login --allow-no-subscriptions --tenant 00000000-0000-0000-0000-000000000000 --scope api://fabric_terraform_provider/default
