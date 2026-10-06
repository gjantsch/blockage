```text
▗▄▄▖ ▗▖    ▗▄▖  ▗▄▄▖▗▖ ▗▖ ▗▄▖  ▗▄▄▖▗▄▄▄▖
▐▌ ▐▌▐▌   ▐▌ ▐▌▐▌   ▐▌▗▞▘▐▌ ▐▌▐▌   ▐▌   
▐▛▀▚▖▐▌   ▐▌ ▐▌▐▌   ▐▛▚▖ ▐▛▀▜▌▐▌▝▜▌▐▛▀▀▘
▐▙▄▞▘▐▙▄▄▖▝▚▄▞▘▝▚▄▄▖▐▌ ▐▌▐▌ ▐▌▝▚▄▞▘▐▙▄▄▖
```
                                                                                
A puzzle game for the terminal, written in pure Go. No graphics library, no TUI framework: just ASCII characters and ANSI escape codes drawn directly to your shell.

The motivation is simple: nostalgia. 

There is something deeply satisfying about ASCII animations on a blank terminal screen, crafted entirely out of whitespace, asterisks, and a handful of control sequences. `blockage` is an exercise in that tradition.

The similarity with any known game is not coincidence. The matrix size was chosen to be ported to an Arduino with a  8x32 led matrix.

<img width="160" height="657" alt="blockage" src="https://github.com/user-attachments/assets/dadd761b-32fa-407e-9e03-22d4919e8682" />

---

## Install

**Requires Go 1.21+.**

```bash
git clone https://github.com/gjantsch/blockage.git
cd blockage
go install .
```

Or run without installing:

```bash
go run .
```

The game must run in a real terminal (not piped). It will exit immediately if stdout is not connected to a TTY.

---

## Controls

| Key | Action |
|---|---|
| `←` `→` | Move block left / right |
| `↑` | Hard drop (instant fall) |
| `r` | Rotate block 90° clockwise |
| `s` | Cycle to the next block shape |
| `e` | Shrink block to a single dot |
| `q`, `Ctrl-C`, `Esc` | Quit |

---

## System design

The codebase is organised around a clean boundary between game logic and rendering.

```
main.go                  game loop, key bindings, board dimensions
internal/
  game/
    block.go             Block struct, shape library, rotation
    matrix.go            board state, collision, row clearing
    renderer.go          Renderer interface
  render/
    terminal.go          ANSI output, raw mode, layout detection
  ui/
    clock/clock.go       spinner animation
    frame/frame.go       ASCII border builder
```

### Terminal layout without a framework

Rather than pulling in `tcell`, `ncurses`, or any TUI library, `blockage` draws its own frame with `fmt.Print` and then issues a **Device Status Report** ANSI query (`\x1b[6n`) to ask the terminal where the cursor ended up. The terminal replies with the row and column, and those coordinates are used to anchor all subsequent board drawing. The board can appear anywhere on screen and the math stays correct.

---

## Dependencies

```
golang.org/x/term   v0.45.0   raw mode, terminal size
golang.org/x/sys    v0.47.0   indirect (required by x/term)
```

No external game or TUI libraries.

---

## License

See [LICENSE](LICENSE).
