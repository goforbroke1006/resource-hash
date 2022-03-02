package job

type Job struct {
	Name string
	Func func() error
}
