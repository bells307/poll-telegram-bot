package poll_options

type FilePollOptionsProvider struct {
	path string
}

func NewFilePollOptions(path string) (FilePollOptionsProvider, error) {
	panic("not implemented")
}

func (p FilePollOptionsProvider) Add(opt *string) error {
	panic("not implemented")
}

func (p FilePollOptionsProvider) Delete(opt *string) error {
	panic("not implemented")
}

func (p FilePollOptionsProvider) List() error {
	panic("not implemented")
}
