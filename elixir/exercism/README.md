# Exercism Elixir — Local Workflow

All terminal, no browser (mostly). Uses `just` to wrap the Exercism CLI and API.

## Dependencies

- [`just`](https://github.com/casey/just) — command runner
- [`exercism`](https://exercism.org/docs/using/solving-exercises/working-locally) — official CLI, must be configured (`exercism configure --token=...`)
- [`jq`](https://jqlang.github.io/jq/) — JSON parsing
- [`bat`](https://github.com/sharkdp/bat) — syntax-highlighted file viewer
- `nvim` — editor

## Typical Workflow

```
just next       # get the next exercise and open it in nvim
just docs       # read the problem description
just hints      # read hints if you're stuck
just test       # run the test suite
just submit     # submit your solution
just complete   # open the exercise page in browser to mark complete (one required click)
just next       # repeat
```

## All Commands

| Command | Description |
|---|---|
| `just next` | Download and open the next unlocked, unsubmitted exercise |
| `just open` | Open the current exercise in nvim |
| `just docs` | Show the exercise instructions |
| `just hints` | Show hints (not all exercises have them) |
| `just test` | Run `mix test` for the current exercise |
| `just submit` | Submit the current exercise to Exercism |
| `just complete` | Open the exercise on exercism.org to click "Mark as complete" |
| `just list` | List all 168 exercises with unlock status |
| `just list-unlocked` | List only unlocked exercises |
| `just list-locked` | List only locked exercises |
| `just get <slug>` | Download a specific exercise and set it as current |
| `just get-all` | Download all currently unlocked exercises |

## How "current" Works

Commands like `open`, `docs`, `test`, `submit`, and `complete` all operate on the
**current exercise**, tracked in a `.current` file at the workspace root.

This file is updated automatically when you run `just next` or `just get <slug>`.
If you ever need to switch manually:

```bash
echo "exercise-slug" > .current
```

## Unlocking Exercises

Exercism concept exercises unlock sequentially — you have to complete one to unlock
the next batch. The unlock flow:

1. `just submit` — submits your code
2. `just complete` — opens the browser so you can click "Mark as complete" (no API for this)
3. `just next` — picks up the newly unlocked exercise
