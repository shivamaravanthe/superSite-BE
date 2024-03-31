package dto

import "unsafe"

type Users struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PasswordStore struct {
	ID          int    `json:"id"`
	Link        string `json:"link"`
	UserName    string `json:"userName"`
	Password    string `json:"password"`
	Description string `json:"description"`
	Master      string `json:"master"`
}

func StringToBytes(s string) []byte {
	p := unsafe.StringData(s)
	b := unsafe.Slice(p, len(s))
	return b
}
