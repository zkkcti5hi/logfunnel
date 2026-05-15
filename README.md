# logfunnel

A structured log aggregator that tails multiple log sources and routes entries to different sinks based on regex filter rules.

---

## Installation

```bash
go install github.com/yourname/logfunnel@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/logfunnel.git && cd logfunnel && go build ./...
```

---

## Usage

Define your sources, filters, and sinks in a config file:

```yaml
sources:
  - name: app
    path: /var/log/app.log

routes:
  - match: "level=error"
    sink: stderr_file
  - match: "level=(info|debug)"
    sink: stdout_file

sinks:
  - name: stderr_file
    path: /var/log/errors.log
  - name: stdout_file
    path: /var/log/general.log
```

Then run:

```bash
logfunnel --config logfunnel.yaml
```

logfunnel will tail all configured sources and route each log entry to the appropriate sink based on the first matching regex rule. Unmatched entries are dropped by default unless a fallback sink is specified.

---

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `logfunnel.yaml` | Path to config file |
| `--dry-run` | `false` | Print routing decisions without writing |
| `--verbose` | `false` | Enable debug output |

---

## License

MIT © yourname