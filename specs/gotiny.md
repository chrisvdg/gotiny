# GoTiny specs

## first release goals

The goal of this project is to make a server where url can be stored and be called through a shorthand name or ID to easily retrieve and/or be redirected to that URL.

The data will be stored in a backend file (YAML/JSON/...) or database (Mongo/MariaDB) if easily feasible. It should at least be modular so one or the other is easily plugged in.

The shorthand should either be custom defined or randomly generated as a short base64 URL safe string. (5 chars?)

Creating shorthand URLS already found in the database will return the ID already in the backend and return a success code.

For authentication, make it possible that creating a new shorthand requires a token.
It should also be possible to only require a token when creating custom shorthand ID's and have different tokens for both options.
Maybe also add a token for read operations (Expand url or list shorthands)
This is provided at startup of the program. The config file, flags or documentation contains the environmental argument where these tokens are stored.

API specs can be found in the [api.yaml](./api.yaml) RAML spec file

TLS support is required.

## Long term features

If there is demand for it, it could be made possible for an admin to create multiple tokens with ACL associated stored in the backend so sharing access can be more controlled and easily revoked to a smaller subset of users. 
Then startup would require an admin token and other tokens would be stored in the backend (hashed?) and API's added to create and revoke these tokens authorized by the admin token.
