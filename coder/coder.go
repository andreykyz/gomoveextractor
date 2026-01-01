package coder

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/cockroachdb/errors"
	"gocv.io/x/gocv"
)

const MinimumArea = 12000

type CoderArgs struct {
	InputFiles       []string
	InputBasePath    *string
	OutputVideoDir   string
	OutputRectangles string
}

type coder = CoderArgs

type Coder interface {
	Generate() error
}

func NewCoder(args CoderArgs) Coder {
	return &args
}

func (c *coder) Generate() error {
	for _, file := range c.InputFiles {
		video, err := gocv.VideoCaptureFile(file)
		if err != nil {
			return errors.Wrapf(err, "Error opening video capture file: %s\nddd%vddd\n", file)
		}
		defer video.Close()

		window := gocv.NewWindow("Motion Window")
		defer window.Close()

		img := gocv.NewMat()
		defer img.Close()

		imgDelta := gocv.NewMat()
		defer imgDelta.Close()

		imgThresh := gocv.NewMat()
		defer imgThresh.Close()

		mog2 := gocv.NewBackgroundSubtractorMOG2()
		defer mog2.Close()

		status := "Ready"

		fmt.Printf("Start reading file: %v\n", file)
		for {
			time.Sleep(1 * time.Second)
			if ok := video.Read(&img); !ok {
				return errors.Newf("File closed: %v\n", file)
			}
			if img.Empty() {
				continue
			}

			status = "Ready"
			statusColor := color.RGBA{0, 255, 0, 0}

			// first phase of cleaning up image, obtain foreground only
			mog2.Apply(img, &imgDelta)

			// remaining cleanup of the image to use for finding contours.
			// first use threshold
			gocv.Threshold(imgDelta, &imgThresh, 15, 255, gocv.ThresholdBinary)

			// then dilate
			kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
			gocv.Dilate(imgThresh, &imgThresh, kernel)
			kernel.Close()

			// now find contours
			contours := gocv.FindContours(imgThresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)
			contSize := contours.Size()
			var areaDetected float64
			for i := 0; i < contSize; i++ {
				area := gocv.ContourArea(contours.At(i))
				if area < MinimumArea {
					continue
				}
				areaDetected = area
				status = "Motion detected"
				statusColor = color.RGBA{255, 0, 0, 0}
				//	gocv.DrawContours(&img, contours, i, statusColor, 2)

				rect := gocv.BoundingRect(contours.At(i))
				gocv.Rectangle(&img, rect, color.RGBA{0, 0, 255, 0}, 2)
			}

			contours.Close()

			gocv.PutText(&img, fmt.Sprintf("%s area - %f", status, areaDetected), image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, statusColor, 2)

			window.IMShow(img)
			if window.WaitKey(10) == 27 {
				break
			}
		}
	}
	return nil
}
