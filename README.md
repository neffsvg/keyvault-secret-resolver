# keyvault-secret-resolver

Resolves secrets, based on their name from a given azure keyvault and writes them into an .env file.

## Installation
TODO - not yet available

## How it works

First add a Key/Value Pair to .env.secrets file.
```dotenv
# .env.secrets (or whatever name you choose)
ENV_VARIABLE=secret-name-in-keyvault
```
- `ENV_VARIABLE` is the key/name of the env variable, it'll be the key/name in the generated output file.
- `secret-name-in-keyvault` is the name of the secret in the keyvault, which will be resolved to its latest value.

After running it, the value of the secret named `secret-name-in-keyvault` will be loaded from azure keyvault and written to an output file.
```dotenv
# .env (or whatever name you choose)
ENV_VARIABLE=some value from keyvault
```

## usage

```
-h, --help                            lists available options
-k, --keyvault              string    keyvault name to get env variables from 
-r, --result-env-file-path  string   .env file to put in config and resolved secrets (default ".env")
-s, --secrets-env-file      string    secrets.env file where secrets are specified with their keyvault names
```

## run the project

```zsh
make help # targets for this project
make build # builds the project
```

from code

```zsh
go run .
go run . -h
go run . -s .env.secrets -k <keyvaultname>
```

or run the build artifact

```zsh
make build
./tmp/keyvault-secret-resolver
./tmp/keyvault-secret-resolver -h
./tmp/keyvault-secret-resolver -s .env.secrets -k <keyvaultname>
```