package contracts

type FileProcessor[T any] func(filepath string) (T, error)

type FilesProcessor[T any] interface {
	ProcessFiles(files []string) (map[string]T, error)
}

type FilesProcessorConfig[T any] struct {
	Processor      FileProcessor[T]
	Recursive      bool
	FallbackGoList bool
}
