package domain

type OptFunc func(*Opts)

type Opts struct {
	repository Repository
}

func defaultOpts(Repository Repository) Opts {
	return Opts{
		repository: Repository,
	}
}
