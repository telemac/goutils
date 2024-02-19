package cli

type Username struct {
	Username string `help:"username" short:"u" json:"username"`
}

type Password struct {
	Password string `help:"password" short:"p" json:"password"`
}
