happydoc
---
convenient system that easy to create to make people easy to publish their documents

## config file - .happydoc.json

```json
{
    "project": "myAwesomeProject",
    "server": "http://127.0.0.1:8000"
}
```

## Commands

- `happydoc init`

generate config file `.happydoc.json` by answer questions. If project has file `package.json`, project will default to name field inside it.

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
