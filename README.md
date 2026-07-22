# TapeRunner

**TapeRunner** (also branded *Escape From The Tape*) is a deterministic single-tape Turing machine simulator written in Go. You can define machines in YAML, run them from a web UI, and also run a single-tape universal Turing machine (UTM) that simulates other machines from an encoded `machine#input` description.

## Features

- Deterministic single-tape TM simulation (read / write / move left or right)
- Machine definitions as YAML files under `machines/`
- Browser UI to pick a machine, enter an input tape, and step through the run
- Universal TM runner (`cmd/utm`) that loads a transition table and interprets encoded machines
- Optional symbol-map decoding so UTM results can be shown in the original alphabet

## Requirements

- [Go](https://go.dev/) 1.26+ (see `go.mod`)

```bash
go mod tidy
```

## Project layout

| Path | Role |
|------|------|
| `cmd/server` | HTTP server + web UI |
| `cmd/utm` | CLI for the universal Turing machine |
| `internal/domain` | Tape, machine, simulator core |
| `internal/application` | Simulation service |
| `internal/adapters/yaml_parser` | Load / validate YAML machines |
| `internal/adapters/http` | Routes and handlers |
| `internal/machines` | UTM table loading, tape encoding, decoding |
| `machines/` | Sample YAML machines and UTM demos |
| `web/` | UI templates and static JS |

## How to use

### 1. Web simulator

Start the server from the repo root (so it can find `machines/` and `web/`):

```bash
go run ./cmd/server
```

Open [http://localhost:3225](http://localhost:3225).

1. Choose a machine from the dropdown (any `*.yaml` in `machines/`).
2. Enter the initial tape (e.g. `1001`).
3. Run the simulation and inspect each step (state, head, tape, status).

Sample machines included:

- `even_palindrome.yaml` — accepts even-length palindromes over `{0,1}`
- `final_gate.yaml` — sample “final gate” machine

### 2. Define your own machine

Add a YAML file under `machines/`. Minimal shape:

```yaml
states:
  - q0
  - accept
  - reject

input_alphabet:
  - "0"
  - "1"

tape_alphabet:
  - "0"
  - "1"
  - "_"

blank: "_"
start_state: q0

accept_states:
  - accept

reject_states:
  - reject

transitions:
  - current_state: q0
    read: "0"
    write: "0"
    move: R
    next_state: q0
  - current_state: q0
    read: "_"
    write: "_"
    move: R
    next_state: accept
```

Restart or refresh the UI; the new file appears in the machine list automatically.

API endpoints used by the UI:

- `GET /api/machines` — list YAML machines
- `POST /api/run` — body `{"machine_name":"even_palindrome.yaml","initial_tape":"1001"}`

### 3. Universal Turing machine (CLI)

The UTM is itself a large TM transition table (`machines/utm_table.tsv`). It takes an encoding of another machine `M` concatenated with `#` and `M`’s input:

```text
rules#input
```

Each rule of `M` is:

```text
current,read,next,write,move;
```

Conventions used by this UTM encoding:

- Start state of `M` is labeled `0`
- Accept / reject of `M` are `A` / `X`
- Tape symbols of `M` are digits `0`–`9` (`0` is blank)
- Input file format: all rules, then `#`, then the input string

Run a demo:

```bash
go run ./cmd/utm machines/utm_table.tsv machines/utm_demo1.txt
```

Optionally pass a symbol map to decode the final tape of `M` into its real alphabet:

```bash
go run ./cmd/utm \
  machines/utm_table.tsv \
  machines/utm_demo5_final_gate_1bit.txt \
  machines/utm_demo5_final_gate.map
```

Example compact encoding (`machines/utm_demo1.txt`):

```text
0,1,0,1,R;0,2,X,2,R;0,0,A,0,R;#111
```

Meaning: machine definition, then `#`, then input `111`.

## Simulation notes

- The ordinary simulator (`domain.Simulator`) stops after **500** steps (`MaxSteps`).
- The UTM runner uses a much higher step budget (default **500000**) because each step of `M` costs many UTM steps.
- Missing transitions cause reject; reaching an accept/reject state ends the run.

## License

MIT — see [LICENSE](LICENSE).
