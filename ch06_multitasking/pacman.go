package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strings"
	"time"
)

const (
	GameWidth  = 20
	GameHeight = 20

	// Delay is application pacing (game logic cadence), not scheduler time slice.
	Delay = 1 * time.Second
)

type point struct{ x, y int }

// define the initial state of the game world
var (
	pacmanPos   = point{0, 0}
	ghosts      = []point{{5, 5}, {10, 10}}
	score       = -10
	isGameOver  = false
	gameOverMsg = ""
	dots        = make(map[point]struct{}, GameWidth*GameHeight)
)

func init() {
	// initialize dots in all positions
	for x := range GameWidth {
		for y := range GameHeight {
			dots[point{x, y}] = struct{}{}
		}
	}
}

// --------------------
// Utility functions
// --------------------

// inBounds checks if a point is within the game boundaries.
func inBounds(p point) bool {
	return 0 <= p.x && p.x < GameWidth && 0 <= p.y && p.y < GameHeight
}

// clampToBounds adjusts a point to be within the game boundaries.
func clampToBounds(p *point) {
	if p.x < 0 {
		p.x = 0
	}
	if p.x >= GameWidth {
		p.x = GameWidth - 1
	}
	if p.y < 0 {
		p.y = 0
	}
	if p.y >= GameHeight {
		p.y = GameHeight - 1
	}
}

// clearScreen clears the terminal screen.
func clearScreen() {
	fmt.Print("\033[2J\033[H") // ANSI clear + cursor home
}

// isGhosts checks if a point is occupied by a ghost.
func isGhosts(p point) bool {
	return slices.Contains(ghosts, p)
}

// --------------------
// Input abstraction
// --------------------

// InputSource defines an interface for retrieving input commands.
type InputSource interface {
	// Get returns a command string and a boolean indicating whether input was available.
	// The block flag tells the source whether to block waiting for input.
	Get(block bool) (string, bool)
}

// BlockingStdinSource blocks on stdin reads (line-buffered; user presses Enter).
//
// This source models the naive, blocking input thread from the non-multitasking version.
type BlockingStdinSource struct {
	r *bufio.Reader
}

func NewBlockingStdinSource() *BlockingStdinSource {
	return &BlockingStdinSource{r: bufio.NewReader(os.Stdin)}
}

func (s *BlockingStdinSource) Get(block bool) (string, bool) {
	// for this source, block is ignored; it's always blocking
	line, err := s.r.ReadString('\n')
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(line), true
}

// SharedBufferSource reads a shared string buffer. It represents non-blocking input polling.
//
// A separate goroutine is expected to fill the buffer asynchronously.
// The multitasking version implements that goroutine function with `readInputBlockingInto`.
type SharedBufferSource struct {
	buf *string
}

func NewSharedBufferSource(buf *string) *SharedBufferSource {
	return &SharedBufferSource{buf: buf}
}

func (s *SharedBufferSource) Get(block bool) (string, bool) {
	// non-blocking by design
	if s.buf == nil || *s.buf == "" {
		return "", false
	}
	cmd := *s.buf
	*s.buf = "" // consume and reset
	return cmd, true
}

// --------------------
// Tasks
// --------------------

type StepTask interface {
	Name() string
	Period() time.Duration // 0 means "every scheduler slice"
	Step(now time.Time, block bool)
}

type InputTask struct {
	src InputSource
}

func (t *InputTask) Name() string          { return "getUserInput" }
func (t *InputTask) Period() time.Duration { return 0 } // can run every slice
func (t *InputTask) Step(now time.Time, block bool) {
	if isGameOver {
		return
	}

	cmd, ok := t.src.Get(block)
	if !ok {
		// in mt mode (block=false), "no input" is normal
		// in no-mt mode (block=true), ok=false can mean stdin closed/error
		if block {
			isGameOver = true
			gameOverMsg = "input error"
		}
		return
	}

	switch cmd {
	case "q":
		isGameOver = true
		gameOverMsg = "quit"
		return
	case "w":
		pacmanPos = point{pacmanPos.x, pacmanPos.y - 1}
	case "a":
		pacmanPos = point{pacmanPos.x - 1, pacmanPos.y}
	case "s":
		pacmanPos = point{pacmanPos.x, pacmanPos.y + 1}
	case "d":
		pacmanPos = point{pacmanPos.x + 1, pacmanPos.y}
	default:
		// ignore
	}
	clampToBounds(&pacmanPos)
}

type WorldTask struct{}

func (t *WorldTask) Name() string          { return "computeGameWorld" }
func (t *WorldTask) Period() time.Duration { return Delay }
func (t *WorldTask) Step(now time.Time, block bool) {
	if isGameOver {
		return
	}

	// move ghosts randomly
	for i, g := range ghosts {
		nx := g.x + []int{-1, 0, 1}[rand.Intn(3)]
		ny := g.y + []int{-1, 0, 1}[rand.Intn(3)]
		np := point{nx, ny}
		if inBounds(np) {
			ghosts[i] = np
		}
	}

	// check collision pacman with ghost
	if slices.Contains(ghosts, pacmanPos) {
		isGameOver = true
		gameOverMsg = "caught"
		return
	}

	// pacman eat dot
	if _, ok := dots[pacmanPos]; ok {
		delete(dots, pacmanPos)
		score += 10
	}

	// win condition
	if len(dots) == 0 {
		isGameOver = true
		gameOverMsg = "win"
		return
	}
}

type RenderTask struct{}

func (t *RenderTask) Name() string          { return "renderNextScreen" }
func (t *RenderTask) Period() time.Duration { return Delay }
func (t *RenderTask) Step(now time.Time, block bool) {
	clearScreen()

	if isGameOver {
		fmt.Println("GAME OVER!")
		fmt.Printf("Your score: %d. (reason: %s)\n", score, gameOverMsg)
		fmt.Println("Press Enter to exit.")
		return
	}

	fmt.Printf("Score: %d. Press 'q' then Enter to quit.\n", score)

	for y := range GameHeight {
		var b strings.Builder
		for x := range GameWidth {
			p := point{x, y}
			char := " "
			if p == pacmanPos {
				char = "P"
			} else if isGhosts(p) {
				char = "G"
			} else if _, ok := dots[p]; ok {
				char = "."
			}
			if x > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(char)
		}
		fmt.Println(b.String())
	}
}
