package files

// Image interface is the one that defines what an images is
// in this context. It's responsible to return kind of the underlying image.
type Image interface {
	ImageType() string
}
