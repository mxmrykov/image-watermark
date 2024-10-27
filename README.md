# Image watermark
Image watermark is a golang app, that allows you to write some ASCII text on input image with basic parameters.
Now it supports only 27 characters - upper cased english letters and a space.
Principle of work - it transforms input text to byte array and get those letters explanation in relation matrix`s of bool data.
Then it parses input image as two-dimensional array of pixels, and rewrite it with added matrix of text.

## Basic usage:

`git clone https://github.com/mxmrykov/image-watermark.git`

After cloning app, open it and to call writer:

```
app, err := internal.NewApp(uint8(marginTop), uint8(marginLeft), uint8(fontSize))
if err != nil {
  log.Fatalln(err)
}

if err = app.WriteText(text); err != nil {
  log.Fatalln(err)
}
```

Replace `marginTop, marginLeft, fontSize` with needed values and put some PNG image into `media/input` folder.

Then, start a program, and in case everything is ok, processed image will located in `media/output`.
