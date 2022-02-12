package poll_options

type (
	PollOptionsProvider interface {
		Add(opt string) error
		Delete(opt string) error
		List() ([]string, error)
	}
)
