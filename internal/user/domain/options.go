package domain

type OptFunc func(*Opts)

type Opts struct {
	userRepository  UserRepository
	authService     AuthService
	passwordService PasswordService
}

func defaultOpts(userRepository UserRepository, authService AuthService, passwordService PasswordService) Opts {
	return Opts{
		userRepository:  userRepository,
		authService:     authService,
		passwordService: passwordService,
	}
}
