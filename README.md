## GO-8 - A CHIP-8 emulator written in Go.

This emulator is a WIP, so while all functionality is implemented, there may be bugs or inconsistencies!

### Installation
Note: requires Go version >= 1.11.

```bash
go get https://github.com/nginth/go-8
```

or

```bash
git clone https://github.com/nginth/go-8
cd go-8
go build
```

### Usage

```
Usage of ./go-8:
  -clockFreq int
    	Clock speed in Hz. (default 300)
  -rom string
    	Path to rom. (default "roms/tetris.ch8")
  -timerFreq int
    	Timer frequency in Hz. (default 60)
```

### ROM Compatibility

This emulator is known to work with the following ROMS:

* [Tetris](https://github.com/dmatlack/chip8/blob/master/roms/games/Tetris%20%5BFran%20Dachille%2C%201991%5D.ch8)
* [Breakout](https://github.com/badlogic/chip8/blob/master/roms/breakout.rom)
