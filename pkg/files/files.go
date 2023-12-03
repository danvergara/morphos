package files

type File interface {
	SupportedFormats() map[string]string
	ConvertTo(string) ([]byte, error)
}
