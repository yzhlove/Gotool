package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type RectangleWithText struct {
	widget.BaseWidget
	textStr     string
	fillColor   color.Color
	strokeColor color.Color
	textColor   color.Color

	rectangle *canvas.Rectangle
	txt       *canvas.Text
}

type RectangleFunc func(c *RectangleWithText)

func WithText(text string) RectangleFunc {
	return func(c *RectangleWithText) {
		c.textStr = text
	}
}

func WithFillColor(color2 color.Color) RectangleFunc {
	return func(c *RectangleWithText) {
		c.fillColor = color2
	}
}

func WithStrokeColor(color2 color.Color) RectangleFunc {
	return func(c *RectangleWithText) {
		c.strokeColor = color2
	}
}

func WithTextColor(color2 color.Color) RectangleFunc {
	return func(c *RectangleWithText) {
		c.textColor = color2
	}
}

func NewRectangleText(opts ...RectangleFunc) *RectangleWithText {

	c := &RectangleWithText{
		textStr:     "",
		fillColor:   color.White,
		strokeColor: color.Black,
		textColor:   color.Black,
	}

	// 参数：textStr - 显示的文本；fillColor - 方形填充色；strokeColor - 方形边框色；textColor - 文本色。
	for _, opt := range opts {
		opt(c)
	}

	c.rectangle = canvas.NewRectangle(c.fillColor)
	c.rectangle.StrokeWidth = 1
	c.rectangle.StrokeColor = c.strokeColor
	c.txt = canvas.NewText(c.textStr, c.textColor)
	c.txt.TextStyle = fyne.TextStyle{Bold: true} // 可选：加粗文本

	c.ExtendBaseWidget(c)
	return c
}

func (c *RectangleWithText) Update(opts ...RectangleFunc) {
	for _, opt := range opts {
		opt(c)
	}
	c.rectangle.FillColor = c.fillColor
	c.rectangle.StrokeColor = c.strokeColor
	c.txt.Text = c.textStr
	c.txt.Color = c.textColor

	canvas.Refresh(c.rectangle)
	canvas.Refresh(c.txt)
}

func (c *RectangleWithText) CreateRenderer() fyne.WidgetRenderer {
	return &rectangleRenderer{
		rectangle: c.rectangle,
		text:      c.txt,
	}
}

func (c *RectangleWithText) MinSize() fyne.Size {
	return c.CreateRenderer().(fyne.WidgetRenderer).MinSize()
}

type rectangleRenderer struct {
	rectangle *canvas.Rectangle
	text      *canvas.Text
}

func (r *rectangleRenderer) Destroy() {
	// 释放资源（此处无需特殊处理）
}

func (r *rectangleRenderer) Layout(size fyne.Size) {

	r.rectangle.Resize(size)

	// 计算文本居中位置
	textSize := r.text.MinSize()
	halfWidth := (size.Width - textSize.Width) / 2
	halfHeight := (size.Height - textSize.Height) / 2
	r.text.Resize(textSize)
	r.text.Move(fyne.NewPos(halfWidth, halfHeight))
}

func (r *rectangleRenderer) MinSize() fyne.Size {
	textSize := r.text.MinSize()
	return fyne.NewSize(textSize.Width+15, textSize.Height+15)
}

func (r *rectangleRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.rectangle, r.text}
}

func (r *rectangleRenderer) Refresh() {
	canvas.Refresh(r.rectangle)
	canvas.Refresh(r.text)
}
