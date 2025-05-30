func wattermarkTest(inputImage string, outputImage string, wattermarkImageFile string) {
	// Open a test image.
	srcImage, err := imgconv.Open(inputImage)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	markImage, err := imgconv.Open(wattermarkImageFile)
	if err != nil {
		log.Fatalf("failed to open wattermark image: %v", err)
	}
	markImage400 := imgconv.Resize(markImage, &imgconv.ResizeOption{Width: 200})
	// Resize srcImage to width = 800px preserving the aspect ratio.
	dstImage800 := imgconv.Resize(srcImage, &imgconv.ResizeOption{Width: 1200})
	dstImage := imgconv.Watermark(dstImage800, &imgconv.WatermarkOption{Mark: markImage400, Opacity: 120, Random: false,
		Offset: image.Pt(-520, 270)})

	newFile, err := os.Create(outputImage)

	// Write the resulting image
	if err := imgconv.Write(newFile, dstImage, &imgconv.FormatOption{Format: imgconv.JPEG}); err != nil {
		log.Fatalf("failed to write image: %v", err)
	}
	defer newFile.Close()

}



	wtmark := "/home/jgnoguer/uocWksp/knative/tempStorage/richard.jpg"
	outmark := "/home/jgnoguer/uocWksp/knative/tempStorage/richard-watter.jpg"
	wattermarkTest(wtmark, outmark, "perroagua.png")