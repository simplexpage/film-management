package domain

type OptFunc func(*Opts)

type Opts struct {
	repository Repository
}

func defaultOpts(repository Repository) Opts {
	return Opts{
		repository: repository,
	}
}
