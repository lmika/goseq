package canvas

type multiCanvas struct {
	canvases []Canvas
}

// MultiCanvas returns a canvas which forwards all events writen to it to
// all the passed in canvases, similar to an io.MultiWriter
func MultiCanvas(canvases ...Canvas) Canvas {
	return &multiCanvas{canvases}
}

func (mc *multiCanvas) Close() error {
	for _, c := range mc.canvases {
		// TODO: Close all canvases, regardless of success
		if err := c.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (mc *multiCanvas) Line(fx int, fy int, tx int, ty int, stroke StrokeStyle) {
	for _, c := range mc.canvases {
		c.Line(fx, fy, tx, ty, stroke)
	}
}

func (mc *multiCanvas) Rect(x int, y int, w int, h int, stroke StrokeStyle, fill FillStyle) {
	for _, c := range mc.canvases {
		c.Rect(x, y, w, h, stroke, fill)
	}
}

func (mc *multiCanvas) Circle(x int, y int, rad int, stroke StrokeStyle, fill FillStyle) {
	for _, c := range mc.canvases {
		c.Circle(x, y, rad, stroke, fill)
	}
}

func (mc *multiCanvas) Text(left int, bottom int, line string, style FontStyle) {
	for _, c := range mc.canvases {
		c.Text(left, bottom, line, style)
	}
}

func (mc *multiCanvas) Polygon(xs []int, ys []int, closed bool, stroke StrokeStyle, fill FillStyle) {
	for _, c := range mc.canvases {
		c.Polygon(xs, ys, true, stroke, fill)
	}
}

func (mc *multiCanvas) Polyline(xs []int, ys []int, stroke StrokeStyle) {
	for _, c := range mc.canvases {
		c.Polyline(xs, ys, stroke)
	}
}

func (mc *multiCanvas) Path(path string, stroke StrokeStyle, fill FillStyle) {
	for _, c := range mc.canvases {
		c.Path(path, stroke, fill)
	}
}

func (mc *multiCanvas) SetSize(width int, height int) {
	for _, c := range mc.canvases {
		c.SetSize(width, height)
	}
}
