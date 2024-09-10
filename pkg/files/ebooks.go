package files

// Ebooker interface is the one that defines what a ebook is
// in this context. It's responsible to return kind of the underlying ebook.
type Ebooker interface {
	EbookType() string
}
