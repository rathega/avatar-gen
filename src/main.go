// generates square based avatar for given constants

package main

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"
)

const (
	sideLen        = 250   // side length of our image
	frameFactor    = 50    // frameFactor*frameFactor small squares filled with same color
	isColor        = true  // if false black-white if true colors
	isColorEqualWt = false // if false no weightage to colors if true all colors have equal weightage
	isRandomFill   = true  // if false frames filled in order if true random pattern
	threshold      = -1    // -1 denotes filling all frames other values denotes how many frames to fill
	numColors      = 2     // no. of colors used in filling frames
	symmetric      = 2     // 0 - no symmetry; 1 - top bottom symmetric; 2 - left right symmetric
)

type rgb struct {
	r uint8
	g uint8
	b uint8
}

var colors = []rgb{}

var (
	frames       [][]int // used to denote if a frame is filled or not
	frameSideLen int     // calculated as sideLen/frameFactor (side length of one frame)
	count        int     // count for changing color of each frame
	curFrameX    int     // used in non-random fill for optimization
	curFrameY    int     // used in non-random fill for optimization
	isSet        bool    // used in non-random fill for optimization
)

// init all frames as not set
func init() {
	frameSideLen = sideLen / frameFactor
	frames = make([][]int, frameSideLen)
	for i := range frames {
		frames[i] = make([]int, frameSideLen)
		for j := 0; j < frameSideLen; j++ {
			frames[i][j] = -1
		}
	}
}

// returns if all frames are set (used only in non-random fill)
func isAllFramesSet() bool {
	for i := 0; i < frameSideLen; i++ {
		for j := 0; j < frameSideLen; j++ {
			if frames[i][j] == -1 {
				return false
			}
		}
	}
	return true
}

// returns a random frame that is not set to color
func getRandomFrame() (int, int) {
	x := rand.Intn(frameSideLen)
	y := rand.Intn(frameSideLen)
	if frames[x][y] != -1 {
		return getRandomFrame()
	}
	frames[x][y] = 0
	return x * frameFactor, y * frameFactor
}

// returns the next frame to color in order
func getNextFrame() (int, int) {
	x, y := curFrameX, curFrameY
	curFrameX = (curFrameX + 1) % frameSideLen
	if curFrameX == 0 {
		curFrameY = (curFrameY + 1) % frameSideLen
	}
	if curFrameX == 0 && curFrameY == 0 {
		isSet = true
	}
	return x * frameFactor, y * frameFactor
}

// generates random colors based on numColors value
func generateColors() {
	if !isColor {
		// in case of non-color same r,g,b values chosen for different shades of black white
		for i := 0; i < numColors; i++ {
			rgbval := uint8(rand.Intn(256))
			colors = append(colors, rgb{r: rgbval, g: rgbval, b: rgbval})
		}
	} else {
		// in case of color, random values of r,g,b are chosen
		for i := 0; i < numColors; i++ {
			r := uint8(rand.Intn(256))
			g := uint8(rand.Intn(256))
			b := uint8(rand.Intn(256))
			colors = append(colors, rgb{r: r, g: g, b: b})
		}
	}
}

// generates random non-symmetric image if 'symmetric' is 0
func generateNonSymmetricImage(img *image.RGBA) {
	for {
		// in case of non-random fill, isAllFramesSet logic is avoided for optimisation
		if !isRandomFill && (count == threshold || isSet) {
			break
		}
		// in case of random fill, isAllFramesSet is used for filling all frames
		if isRandomFill && (count == threshold || isAllFramesSet()) {
			break
		}

		// get the point x,y
		var x, y int
		if isRandomFill {
			x, y = getRandomFrame()
		} else {
			x, y = getNextFrame()
		}

		// find the color
		var col color.RGBA
		if isColorEqualWt {
			val := colors[count%numColors]
			col = color.RGBA{val.r, val.g, val.b, 0xff}
		} else {
			val := colors[rand.Intn(numColors)]
			col = color.RGBA{val.r, val.g, val.b, 0xff}
		}

		// set the color in the entire frame
		for i := x; i < (x + frameFactor); i++ {
			for j := y; j < (y + frameFactor); j++ {
				img.Set(i, j, col)
			}
		}

		// count for calculating what color to pick next
		count++
	}
}

// sub function to return top-bottom symmetric frame
func getTopBottomSymmetricFrame(x, y int) (int, int) {
	x, y = x/frameFactor, y/frameFactor
	if -1*(0-y) < frameSideLen-y-1 {
		y = frameSideLen - (-1 * (0 - y)) - 1
	} else {
		y = 0 + (frameSideLen - y - 1)
	}
	frames[x][y] = 0
	return x * frameFactor, y * frameFactor
}

// sub function to return left-right symmetric frame
func getLeftRightSymmetricFrame(x, y int) (int, int) {
	x, y = x/frameFactor, y/frameFactor
	if -1*(0-x) < frameSideLen-x-1 {
		x = frameSideLen - (-1 * (0 - x)) - 1
	} else {
		x = 0 + (frameSideLen - x - 1)
	}
	frames[x][y] = 0
	return x * frameFactor, y * frameFactor
}

// generates symmetric image based on 'symmetric' value 1 or 2
func generateSymmetricImage(img *image.RGBA) {
	// frameCount keeps track of no. of frames set for threshold check purpose
	frameCount := 0
	for {
		// for symmetric fill by default assumed as random fill
		// in case of random fill, isAllFramesSet is used for filling all frames
		if frameCount == threshold || isAllFramesSet() {
			break
		}

		// get the point x,y
		var x, y int
		x, y = getRandomFrame()

		// find the color
		var col color.RGBA
		if isColorEqualWt {
			val := colors[count%numColors]
			col = color.RGBA{val.r, val.g, val.b, 0xff}
		} else {
			val := colors[rand.Intn(numColors)]
			col = color.RGBA{val.r, val.g, val.b, 0xff}
		}

		// set the color in the entire frame
		for i := x; i < (x + frameFactor); i++ {
			for j := y; j < (y + frameFactor); j++ {
				img.Set(i, j, col)
			}
		}
		frameCount++

		// count for calculating what color to pick next
		count++

		var x1, y1 int
		// pick the symmetric frame
		if symmetric == 1 {
			x1, y1 = getTopBottomSymmetricFrame(x, y)
		} else {
			x1, y1 = getLeftRightSymmetricFrame(x, y)
		}

		// skip the frame if symmetric frame not present
		// (happens in odd number of frames cases)
		if x == x1 && y == y1 {
			continue
		}

		// set the same color in the symmetric frame
		for i := x1; i < (x1 + frameFactor); i++ {
			for j := y1; j < (y1 + frameFactor); j++ {
				img.Set(i, j, col)
			}
		}
		frameCount++
	}
}

func main() {

	// setting seed as time for different results each time for same constants
	rand.Seed(time.Now().UTC().UnixNano())

	topLeft := image.Point{0, 0}
	bottomRight := image.Point{sideLen, sideLen}

	img := image.NewRGBA(image.Rectangle{topLeft, bottomRight})

	// generate the required no. of colors 'numColors'
	generateColors()

	// call symmetric or non-symmetric gen function based on 'symmetric'
	if symmetric == 1 || symmetric == 2 {
		generateSymmetricImage(img)
	} else {
		generateNonSymmetricImage(img)
	}

	// encode as PNG
	f, _ := os.Create("image.png")
	png.Encode(f, img)
}
