package service

type Session struct {
	Account string
	Password string
}

// Login user login
func (s *Session) Login() {
	session := GetSession()
	defer DeleteSession(session, true)
	// TODO: too many things to do
}