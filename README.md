# run 🚀

> *THE* universal & ecosystem independent lightweight script management tool!

- Are you tired of write a bunch of shell scripts just to execute a single one-liner?
- Do you like javascript - package.json approach to manage your scripts and wish to integrate that style of script management into other ecosystems?

__*`run`*__ is your solution! (It even autocompletes✨)

## Installation

```bash
go install github.com/fabiankachlock/run
```

## Quick Start

1) create a `run.json` in your project root
2) define your scripts

```json
{
    "scripts": {
        "build": "go build",
        "doSomething": "echo 'Hello World'"
    }
}
```

> For further configuration please see [config-guide.md](https://github.com/fabiankachlock/run/blob/main/config-guide.md)

3) use it! `run build`
4) setup auto completion

Put  `source <(run --completion)` into your `.bashrc`, `.zshrc` or according

1) integrate existing script providers

> *run.json*
> ```json
> {
>     "scopes": {
>         "npm": true // makes all your package.json scripts available to run
>     }
> }
> ```

> *package.json*
> ```json
> {
>     "scripts": {
>         "dev": "node ...",
>     }
> }
> ```

`run npm:dev`


## Concepts

### auto completion

To setup auto completion in your shell, you need to put the following line into your shell configuration file (e.g. `.bashrc` or `.zshrc`)

```bash
source <(run --completion)
```

### writing scripts

TODO

### scopes & extended configuration

TODO

### integration with existing scripts

External script providers are just special scopes which need to be enabled explicitly.

```json
{
    "scopes": {
        "npm": true, // scripts are available with "npm:" prefix 
        "yarn": "myYarn", // scripts are available with "myYarn:" prefix
        "<key>": <"alias string" | true>
    }
}
```

Currently only the following are supported:
- npm via `npm run` (key: npm)
- yarn via `yarn run` (key: yarn)

## Roadmap

- [x] shell auto completion support
- [x] support `root` option in config
- [x] provide `--list` option
- [x] support global config
- [ ] provide prebuilt binaries in releases
- [ ] support more vendors