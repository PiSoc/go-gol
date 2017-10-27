package main

import (
	"flag"
	"image/color"
    "time"
    "math/rand"
//    "fmt"
	"github.com/mcuadros/go-rpi-rgb-led-matrix"
)

var (
	rows       = flag.Int("led-rows", 32, "number of rows supported")
	parallel   = flag.Int("led-parallel", 1, "number of daisy-chained panels")
	chain      = flag.Int("led-chain", 1, "number of displays daisy-chained")
	brightness = flag.Int("brightness", 100, "brightness (0-100)")
)

func main() {
    // Matrix setup
	config := &rgbmatrix.DefaultConfig
	config.Rows = *rows
	config.Parallel = *parallel
	config.ChainLength = *chain
	config.Brightness = *brightness

    // GoL variable setup
    board := genRandomBoard()
    neighbors := genBlankBoard()

	m, err := rgbmatrix.NewRGBLedMatrix(config)
	fatal(err)

	c := rgbmatrix.NewCanvas(m)
	defer c.Close()

    for {

        // Draw living cells
	    bounds := c.Bounds()
	    for x := bounds.Min.X; x < bounds.Max.X; x++ {
	        for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
                if board[x][y] == 1 {
                    c.Set(x, y, color.RGBA{50, 50, 255, 150})
                }
		    }
	    }
        c.Render()

        // Count neighbors
	    for x := 0; x < 32; x++ {
		    for y := 0; y < 32 ;y++ {
                neighbors[y][x] = countNeighbors(y ,x , 32, board)
            }
        }

        // Update board
        for x := 0; x < 32; x++{
		    for y := 0; y < 32; y++{
                board[y][x] = nextTick(board[y][x], neighbors[y][x])
            }
        }

        //0.08 sec
        time.Sleep(80 *1000 *1000)

    }

}

func genBlankBoard() [32][32]int {
    size := 32

    var board [32][32]int
    for y := 0; y < size; y++ {
        for x := 0; x < size; x++ {
            board[x][y] = 0
        }
    }

    return board
}


func genRandomBoard() [32][32]int {
    randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

    size := 32

    var board [32][32]int
    for y := 0; y < size; y++ {
        for x := 0; x < size; x++ {
            board[x][y] = randGen.Intn(2)
        }
    }

    return board
}


func countNeighbors(x ,y ,size int, board [32][32]int) int {
    var n int
    n = board[x][mod(y-1, size)] +
        board[mod(x+1, size)][mod(y-1, size)] +
        board[mod(x+1, size)][y] +
        board[mod(x+1, size)][mod(y+1, size)] +
        board[x][mod(y+1, size)] +
        board[mod(x-1, size)][mod(y+1, size) ]+
        board[mod(x-1, size)][y] +
        board[mod(x-1, size)][mod(y-1, size)]
    return n
}

func mod(n , size int) int {
    size -= 1
    if n < 0 {
        return size
    } else if n > size {
        return 0
    }
    return n
}

func nextTick(cell, neighbors int) int {
    if cell == 1 && (neighbors == 2 || neighbors == 3) {
        return 1
    } else if cell == 0 && neighbors == 3 {
        return 1
    }
    return 0
}



func init() {
	flag.Parse()
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
