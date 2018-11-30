happydoc
---
convenient system that easy to create to make people easy to publish their documents

## config file - .happydoc.json

```json
{
    "project": "myAwesomeProject",
    "server": "http://127.0.0.1:8000",
    "account": "your account name at server",
    "token":
}
```

## Commands

- `happydoc init`

generate config file `.happydoc.json` by answer questions. If project has file `package.json`, project will
default to name field inside it.

- `happydoc publish`

upload documents to server. It receives one argument to indicate **path to document folder**, and has four flags:
    - `server`: shorthand `s`, will override value inside `.happydoc.json` if used
    - `project`: shorthand `p`, will override value inside `.happydoc.json` if used
    - `version`: shorthand `v`, current version of documents
    - `type`: shorthand `t`, type of documents

    ```bash
    happydoc publish docs/fontend -v 1.2.0 -t frontend
    happydoc publish docs/backend -v 1.2.0 -t backend -s http://my-doc-server.com -p myAwesomeProject2
    ```

- `happydoc server`

init and run happydoc server. This server will receives documents uploaded by happdoc client with `publish` command.
And `token` used by happdoc client is also generated on it.

This command will ask you:
- folder to init server
- which port should be used by server
- password set for postgresql **running on docker container, NOT DB ON YOUR HOST**

# Build and Release
- release using [goreleaser](https://github.com/goreleaser/goreleaser),
release require github token with repo, check [doc](https://goreleaser.com/quick-start/)

    ```bash
    git tag version
    export GITHUB_TOKEN='Your Github Token Here'
    goreleaser goreleaser release

    # dry run without publish to github etc.
    goreleaser release --skip-publish
    ```
- publish release to npm using [go-npm](https://github.com/sanathkr/go-npm): `npm publish`

## This project use
- [golang](https://github.com/golang/go)
- [postgres](https://www.postgresql.org/)
- [cobra](https://github.com/spf13/cobra)
- [gin](https://github.com/gin-gonic/gin)

## TODO
- better UI
- use https
- fix output for subcommand `server`
