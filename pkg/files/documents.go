package files

// Document interface is the one that defines what a document is
// in this context. It's responsible to return kind of the underlying document.
type Document interface {
	DocumentType() string
}
