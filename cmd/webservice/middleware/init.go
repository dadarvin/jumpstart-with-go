package middleware

type Middleware struct {
	secret string
}

func New(secret string) *Middleware {
	return &Middleware{
		secret: secret,
	}
}
