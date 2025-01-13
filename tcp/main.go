package tcp

import (
	. "mandelbrot/mandelbrot"
)

var mandelbrot Mandelbrot

func main() {
	TcpConnection(mandelbrot)

	return
}
