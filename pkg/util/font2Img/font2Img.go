package font2Img

import (
	"bytes"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"regexp"
)

var f *truetype.Font

func init() {
	fontBytes, err := ioutil.ReadFile("./template/msyh.ttc")
	if err != nil {
		panic(err)
		return
	}
	f, err = truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
		return
	}
}

// GetRGBA 将姓名转换为RGBA
func GetRGBA(name string) *image.RGBA {
	tep := []rune(name)
	tepLen := len(tep)

	if tepLen <= 1 {
		return drawImg(name, 72)
	}

	// 存在字母
	if isLetter(name) {
		return drawImg(string(tep[0]), 72)
	}

	return drawImg(string(tep[tepLen-2:tepLen]), 36)
}

// GetReader 图片转换为Reader(
func GetReader(name string) (*bytes.Reader, error) {
	rgba := GetRGBA(name)

	proverbs := new(bytes.Buffer)
	err := png.Encode(proverbs, rgba)
	if err != nil {
		return nil, err
	}

	readSeeker := bytes.NewReader(proverbs.Bytes())

	return readSeeker, nil
}

// 判断字符串中是否包含字母
func isLetter(s string) bool {
	r := regexp.MustCompile("[a-zA-Z]|\\d| ")
	return r.MatchString(s)
}

// drawImg 开始生成img
func drawImg(text string, size float64) *image.RGBA {
	// 创建背景和画板
	fg, bg := image.White, image.NewUniform(color.RGBA{R: 67, G: 162, B: 255, A: 255})
	const imgW, imgH = 120, 120
	rgba := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	draw.Draw(rgba, rgba.Bounds(), bg, image.Pt(0, 0), draw.Src)

	// 绘制文本。
	dpi := float64(72)
	d := &font.Drawer{
		Dst: rgba,
		Src: fg,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    size,
			DPI:     dpi,
			Hinting: font.HintingNone,
		}),
	}

	d.Dot = fixed.Point26_6{
		X: (fixed.I(imgW) - d.MeasureString(text)) / 2,
		Y: fixed.I(imgH+int(size*0.7)) / 2,
	}
	d.DrawString(text)

	return rgba
}
