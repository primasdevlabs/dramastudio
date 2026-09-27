package media

type MediaProcessor interface {
	Process(inputPath, outputPath string) error
}
