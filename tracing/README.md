# Tracing
This example demonstrates the use case of tracing a Go application.
The application in this example generates the mandelbrot set and writes it to an image.
By default the application is filling the image pixel by pixel. The `mode` flag can be used to switch between different
ways to fill the image. The following modes are available:
- `seq`: Fills the image pixel by pixel.
- `px`: Starts width x height number of goroutines to fill each pixel.
- `row`: Starts goroutine for each row of the image to fill the pixels in the row.
- `workers`: Starts `n` goroutines to fill each pixel individually. To change the number of goroutines the `workers`
can be used to set the value. It defaults to 1.

For more information see https://www.youtube.com/watch?v=ySy3sR1LFCQ on how to use the tracing tool.