package domain

type OptFunc func(*Opts)

type Opts struct {
	filmRepository FilmRepository
}

func defaultOpts(filmRepository FilmRepository) Opts {
	return Opts{
		filmRepository: filmRepository,
	}
}
