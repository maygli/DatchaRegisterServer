package datamodel

type UserProfile struct {
	ID     uint
	UserID uint
	Name   string
	Avatar string
}

func (profile UserProfile) IsValid() bool {
	return profile.Name != ""
}
