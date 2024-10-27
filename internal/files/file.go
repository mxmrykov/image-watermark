package files

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"sync"

	"watermark/internal/symbols"
	"watermark/utils/ASCII"

	"github.com/google/uuid"
)

type IParser interface {
	GetPixels() ([][]Pixel, error)
	WritePixels(img [][]Pixel, text string) error
}

type (
	// Pixel - base pixel description with RGBA parameters
	Pixel struct {
		R int
		G int
		B int
		A int
	}

	Parser struct {
		fileName                          string
		File                              *os.File
		RWM                               *sync.RWMutex
		RelMatrixASCII                    symbols.ASCII_SYM
		fontWeight, marginTop, marginLeft uint8
	}
)

const (
	// INPUT_FILE_PATH - directory where input images have to be located in
	INPUT_FILE_PATH = "./media/input/"
	// OUTPUT_FILE_PATH - directory where ready images are placing
	OUTPUT_FILE_PATH = "./media/output/"
)

// NewFileParser - Constructor of Parser Interface
func NewFileParser(marginTop, marginLeft, fontSize uint8) (IParser, error) {
	imgName := *new(string)
	content, _ := os.ReadDir(INPUT_FILE_PATH)

	for _, file := range content {
		if !file.IsDir() {
			imgName = file.Name()
			break
		}
	}

	in, err := os.Open(INPUT_FILE_PATH + imgName)
	defer func() {
		_ = in.Close()
	}()

	if err != nil {
		return nil, err
	}

	newName := fmt.Sprintf("%s.%s", uuid.New().String(), "png")
	out, err := os.Create(OUTPUT_FILE_PATH + newName)
	defer func() {
		_ = out.Close()
	}()

	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(out, in); err != nil {
		return nil, err
	}

	return &Parser{
		fileName:       OUTPUT_FILE_PATH + newName,
		RWM:            new(sync.RWMutex),
		RelMatrixASCII: symbols.GetASCIIRel(),
		marginTop:      marginTop,
		marginLeft:     marginLeft,
		fontWeight:     fontSize,
	}, nil

}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func (p *Parser) mutexRead(target func(file *os.File) ([][]Pixel, error)) ([][]Pixel, error) {
	err := *new(error)

	p.RWM.RLock()
	defer p.RWM.RUnlock()

	p.File, err = os.Open(p.fileName)

	defer func() {
		_ = p.File.Close()
	}()

	if err != nil {
		return nil, err
	}

	return target(p.File)
}

func (p *Parser) mutexWrite(img [][]Pixel, text string, target func(img [][]Pixel, text string) error) error {
	err := *new(error)

	p.RWM.Lock()
	defer p.RWM.Unlock()

	p.File, err = os.Open(p.fileName)

	defer func() {
		_ = p.File.Close()
	}()

	if err != nil {
		return err
	}

	return target(img, text)
}

func (p *Parser) GetPixels() ([][]Pixel, error) {
	return p.mutexRead(func(file *os.File) ([][]Pixel, error) {
		image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
		img, err := png.Decode(p.File)

		if err != nil {
			return nil, err
		}

		bounds := img.Bounds()
		width, height := bounds.Max.X, bounds.Max.Y

		var pixels [][]Pixel
		for y := range height {
			row := make([]Pixel, 0)

			for x := range width {
				row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
			}

			pixels = append(pixels, row)
		}

		return pixels, nil
	})
}

func (p *Parser) WritePixels(img [][]Pixel, text string) error {
	return p.mutexWrite(img, text, func(img [][]Pixel, text string) error {
		bytes := []byte(text)
		FinalRelationMatrix := make([][]bool, 0, len(bytes)*2-1)

		for idx, byte := range bytes {
			a, ok := p.RelMatrixASCII[byte]
			v := upscaleMatrix(a, int(p.fontWeight))
			if !ok {
				return fmt.Errorf("unknown symbol: %v", string(byte))
			}

			if idx == 0 {
				FinalRelationMatrix = v
			} else {
				for ixInner := range v {
					FinalRelationMatrix[ixInner] = append(FinalRelationMatrix[ixInner], v[ixInner]...)
				}
			}

			if idx != len(bytes)-1 {
				sp := upscaleMatrix(ASCII.SpecSymbols[0].RelationMatrix, int(p.fontWeight))
				for ixInner := range sp {
					FinalRelationMatrix[ixInner] = append(FinalRelationMatrix[ixInner], sp[ixInner]...)
				}
			}
		}

		fImg, err := p.drawRelationMatrix(img,
			func(fmb [][]bool) map[[2]int]bool {
				fmm := make(map[[2]int]bool)
				for y := range fmb {
					for x, v := range fmb[y] {
						fmm[[2]int{x + int(p.marginLeft+p.fontWeight), y + int(p.marginTop+p.fontWeight)}] = v
					}
				}

				return fmm
			}(FinalRelationMatrix),
		)
		if err != nil {
			return nil
		}

		height := len(img)
		width := len(img[0])

		if height == 0 || width == 0 {
			return fmt.Errorf("one of dimensions of image is null")
		}

		targetImg := image.NewRGBA(image.Rect(0, 0, width, height))

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				p := fImg[y][x]
				targetImg.Set(x, y, color.RGBA{
					R: uint8(p.R),
					G: uint8(p.G),
					B: uint8(p.B),
					A: uint8(p.A),
				})
			}
		}

		file, err := os.Create(p.fileName)
		if err != nil {
			return err
		}
		defer func() {
			_ = file.Close()
		}()

		return png.Encode(file, targetImg)
	})
}

func (p *Parser) drawRelationMatrix(img [][]Pixel, fm map[[2]int]bool) ([][]Pixel, error) {
	dist := make([][]Pixel, len(img))

	for y := range img {
		dist[y] = make([]Pixel, len(img[y]))

		for x := range img[y] {
			if v, ok := fm[[2]int{x, y}]; !ok || !v {
				dist[y][x].R = img[y][x].R
				dist[y][x].G = img[y][x].G
				dist[y][x].B = img[y][x].B
				dist[y][x].A = img[y][x].A
				continue
			}

			dist[y][x].R = 0
			dist[y][x].G = 0
			dist[y][x].B = 0
			dist[y][x].A = 1
		}
	}

	return dist, nil
}

func upscaleMatrix(arr [][]bool, scale int) [][]bool {
	rows := len(arr)
	cols := len(arr[0])

	newArr := make([][]bool, rows*scale)
	for i := range newArr {
		newArr[i] = make([]bool, cols*scale)
	}

	for i := range rows {
		for j := range cols {
			for x := range scale {
				for y := range scale {
					newArr[i*scale+x][j*scale+y] = arr[i][j]
				}
			}
		}
	}

	return newArr
}
