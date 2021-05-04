package yaas

type Domain struct {
	ID     int    `json:"id"`
	UserID int    `json:"-"`
	Code   string `json:"code"`
	Domain string `json:"name"`
	Desc   string `json:"description"`
}
