package service

type Session struct {
	Account string
	Password string
}

// Login user login
func (s *Session) Login() {}