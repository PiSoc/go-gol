package main

import (
	"flag"
	"image/color"
    "time"
    "math/rand"
//    "fmt"
	"github.com/mcuadros/go-rpi-rgb-led-matrix"
)
//TODO - matrix in a 1D array, rgb cycle
var (
	rows       = flag.Int("led-rows", 32, "number of rows supported")
	parallel   = flag.Int("led-parallel", 1, "number of daisy-chained panels")
	chain      = flag.Int("led-chain", 1, "number of displays daisy-chained")
	brightness = flag.Int("brightness", 100, "brightness (0-100)")
)
    boardWidth := 32

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
		i := 1
		if i>24 {
			i = 0
		}	
		R := 128+sin((i*3+0)*1.3)*128;
		G := 128+sin((i*3+1)*1.3)*128;
		B := 128+sin((i*3+2)*1.3)*128;
		alpha := 150
		
        // Draw living cells
	    bounds := c.Bounds()
	    for x := bounds.Min.X; x < bounds.Max.X; x++ {
	        for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
                if getInt(board,x,y) == 1 {
                    c.Set(x, y, color.RGBA{R, G, B, alpha})
                }
		    }
	    }
        c.Render()

        // Count neighbors
	    for x := 0; x < 32; x++ {
		    for y := 0; y < 32 ;y++ {
                getIntPoint(neighbors,y,x) = countNeighbors(y ,x , &board)
            }
        }

        // Update board
        for x := 0; x < 32; x++{
		    for y := 0; y < 32; y++{
                getIntPoint(board,y,x) = nextTick(getInt(board,y,x),getInt(neighbors,y,x))
            }
        }

        //0.08 sec
        time.Sleep(80 *1000 *1000)
		i++
    }

}
//get value at a location in an array
func getInt(arr []int, x, y int) int{
    return arr[(y*boardWidth)+x]
}
//get pointer at a value in an array
func getIntPoint(arr []int, x, y int) *int{
    return *arr[(y*boardWidth)+x]
}

//get blank board
func genBlankBoard() *[]int{
    size := &boardWidth
    board := make([]int, size*size)
    return *board
}

//get random board
func genRandomBoard() *[]int {
    randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

    size := &boardWidth

    board:= make([]int, size*size)
    for y := 0; y < size; y++ {
        for x := 0; x < size; x++ {
            getIntPoint(board,x,y) = randGen.Intn(2)
        }
    }

    return *board
}

func countNeighbors(x ,y , *board) int {
    var n int
    size := boardWidth
    n = getInt(board,x,mod(y-1, size)) +
        getInt(board,mod(x+1, size),mod(y-1, size)) +
        getInt(board,mod(x+1, size),y) +
        getInt(board,mod(x+1, size),mod(y+1, size)) +
        getInt(board,x,mod(y+1, size)) +
        getInt(board,mod(x-1, size),mod(y+1, size))+
        getInt(board,mod(x-1, size),y) +
        getInt(board,mod(x-1, size),mod(y-1, size))
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
//TODO
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
