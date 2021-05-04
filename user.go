package yaas

type User struct {
	ID           int    `json:"-"`
	Username     string `json:"username"`
	PasswordHash []byte `json:"-"`
	Salt         []byte `json:"-"`
}
