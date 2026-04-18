package interfaces

type ImageProvider interface {
	SearchImage(query string) (string, error)
}
