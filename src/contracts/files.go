package contracts

type FileProcessor[T any] func(filepath string) (T, error)

type FilesProcessor[T any] interface {
	ProcessFiles(files []string) (map[string]T, error)
}
