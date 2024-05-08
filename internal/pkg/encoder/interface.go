package encoder

type VideoEncoder interface {
	Download(filePath string, videoID string) error
	Fragment(videoID string) error
	Encode(videoID string) error
	Upload(videoID string) error
	CleanupFiles(videoID string) error
}
