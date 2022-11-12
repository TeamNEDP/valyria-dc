package model

type UserInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Rating int    `json:"rating"`
}

func (u User) Info() UserInfo {
	return UserInfo{
		ID:     u.ID,
		Name:   u.Name,
		Email:  u.Email,
		Rating: u.Rating,
	}
}

type ScriptInfo struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (s UserScript) Info() ScriptInfo {
	return ScriptInfo{
		Name: s.Name,
		Code: s.Code,
	}
}
