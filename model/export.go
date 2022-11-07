package model

type UserInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u User) Info() UserInfo {
	return UserInfo{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
