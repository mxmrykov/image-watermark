package internal

import (
	"fmt"
	"time"

	"watermark/internal/files"
)

type IApp interface {
	WriteText(text string) error
}

type App struct {
	FP files.IParser
}

func NewApp(marginTop, marginLeft, fontSize uint8) (IApp, error) {
	FP, err := files.NewFileParser(marginTop, marginLeft, fontSize)
	if err != nil {
		return nil, err
	}
	return &App{
		FP: FP,
	}, nil
}

func (a *App) WriteText(text string) error {
	pxls, err := a.FP.GetPixels()
	if err != nil {
		return err
	}

	fmt.Printf("Height: %dpx, Width: %dpx, Text: %s\n", len(pxls), len(pxls[0]), text)
	fmt.Println("Generating text...")
	t := time.Now()

	if err = a.FP.WritePixels(pxls, text); err != nil {
		return err
	}

	fmt.Printf("Generation is over, time spent: %vms\n", time.Since(t).Milliseconds())

	return nil
}
