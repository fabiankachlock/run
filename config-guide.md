# Configuration Guide

General structure:
```json
{
    "scripts": {
        "myScript": "echo 'Hello World'"
    },
    "extends": ["./nested/folder/exec.json"],
    "scopes": {
        "self": "frontend",
        "npm": true,
        "yarn": "otherNpm"
    },
}
```

## `scripts`

Define custom scripts to run when their alias is called.

Syntax: ` [<alias>]: <script>`

## `extends`

Reference other `exec.json` config files to search when calling a script.

> Note: since exec searches the tree upwards for config files with a matching definition it's enough when only the "root" config file references all other configs
> Example: 
> ```
> | -- project
>      | -- exec.json // (1) <- root config file
>      | -- folderA
>           | -- exec.json // (2)
>      | -- folderB
>           | -- exec.json // (3)
> ``` 
> Only the root exec.json (1) needs to reference (2) & (3) to make all scripts available in the entire project tree 

## `scopes`

This section can be used for renaming scopes or enabling vendor loaders.

### Renaming the own scope

> Note: the own scope is disabled by default, but can be enforced with this setting

setting `"$self": true` will enforce scoped access to all scripts configure in this file. The default scope used is the directory name.

Using another name instead of the directory name can be done by setting `"$self": "<alias to use>"`

### Enabling vendor loaders

exec can be integrated into existing developer configurations using vendor loaders. For example by enabling the `npm` loader all scripts defined in the package.json in teh current directory are available with the `npm` prefix.

Supported loaders:
- npm
- yarn

Loaders need to be enabled explicitly by setting `"<loader alias>": true` or can be renamed by setting `"<loader alias>: "<new alias>`