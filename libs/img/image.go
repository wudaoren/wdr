package img

/*
作者：悟道人
2019-06-13 22:00
*/
import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"os"
	"path"
	"strings"
)

const (
	WATER_UPPER_LEFT  int8 = 1 //左上
	WATER_UPPER_RIGHT      = 2 //右上
	WATER_LOWER_LEFT       = 3 //左下
	WATER_LOWER_RIGHT      = 4 //右下
	WATER_FULL             = 5 //平铺
)

//图片处理
type Image struct {
}

//
func New() *Image {
	return new(Image)
}

//打开图片
func (this *Image) Open(fileSrc string) (image.Image, error) {
	file, err := os.Open(fileSrc)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return this.Decode(file, path.Ext(fileSrc))
}

//读取图片
func (this *Image) Decode(f io.Reader, ext string) (image.Image, error) {
	var (
		img image.Image
		err error
	)
	ext = strings.ToLower(ext)
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(f)
		if err != nil {
			return nil, err
		}
	//
	case ".png":
		img, err = png.Decode(f)
		if err != nil {
			return nil, err
		}
	//
	default:
		return nil, errors.New("不支持的类型")
	}
	return img, nil
}

//压缩至指定大小(w和h二选一)
//w=图片宽
//h=图片高
func (this *Image) Resize(img image.Image, w, h int) (image.Image, error) {
	r := img.Bounds()                              //原图大小
	var rate = float64(r.Max.X) / float64(r.Max.Y) //缩放比例
	if w > 0 {
		h = int(float64(w) / rate)
	} else if h > 0 {
		w = int(float64(h) * rate)
	} else {
		return img, nil
	}
	curw, curh := r.Dx(), r.Dy()
	imgNew := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// Get a source pixel.
			subx := x * curw / w
			suby := y * curh / h
			r32, g32, b32, a32 := img.At(subx, suby).RGBA()
			r := uint8(r32 >> 8)
			g := uint8(g32 >> 8)
			b := uint8(b32 >> 8)
			a := uint8(a32 >> 8)
			imgNew.SetRGBA(x, y, color.RGBA{r, g, b, a})
		}
	}
	return imgNew, nil
}

//添加水印，为了美观，只支持图片水印
//file 水印图片路径，
//position 水印坐标 1 左上角，2右上角，3左下角，4右下角，5平铺
func (this *Image) Watermark(img, water image.Image, position int8) (image.Image, error) {
	imgBounds := img.Bounds()     //原始图片网格
	waterBounds := water.Bounds() //水印图片网格
	wx, wy := waterBounds.Max.X, waterBounds.Max.Y
	ix, iy := imgBounds.Max.X, imgBounds.Max.Y
	var rects = make([]image.Rectangle, 0)
	switch position {
	case WATER_UPPER_LEFT: //上
		rects = append(rects, image.Rect(0, 0, wx, wy))
	case WATER_UPPER_RIGHT: //右
		rects = append(rects, image.Rect(ix-wx, 0, ix, wy))
	case WATER_LOWER_LEFT: //下
		rects = append(rects, image.Rect(ix-wx, iy-wy, ix, iy))
	case WATER_LOWER_RIGHT: //左
		rects = append(rects, image.Rect(0, iy-wy, wx, iy))
	case WATER_FULL: //平铺
		x_num := int(math.Ceil(float64(ix) / float64(wx)))
		y_num := int(math.Ceil(float64(iy) / float64(wy)))
		for n := 0; n < y_num; n++ {
			for r := 0; r < x_num; r++ {
				rects = append(rects, image.Rect(r*wx, n*wy, (r+1)*wx, (n+1)*wy))
			}
		}
	}
	baseImage := image.NewNRGBA(imgBounds)
	draw.Draw(baseImage, imgBounds, img, image.ZP, draw.Src)
	for _, rect := range rects {
		draw.Draw(baseImage, rect, water, image.Pt(0, 0), draw.Over)
	}
	return baseImage, nil
}

//根据文件名类型压缩
func (this *Image) Encode(img image.Image, ext string) (f *bytes.Buffer, err error) {
	f = new(bytes.Buffer)
	ext = strings.ToLower(ext) //获取文件后缀
	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(f, img, &jpeg.Options{90}) //90%的压缩比
	case ".png":
		err = png.Encode(f, img)
	default:
		err = errors.New("图片类型错误")
	}
	return
}
