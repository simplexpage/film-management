package domain

// Middleware is a Service type for chainable behavior modifier.
type Middleware func(Service) Service
