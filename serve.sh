#!/usr/bin/env bash

cd dev_root || exit 1

export SHELL_VAULT_ENV=DEV
export SHELL_VAULT_KEK='testkekkey'
export SHELL_VAULT_ROOT_PASSWORD='testtest'
../bin/server -c etc/shell_vault
