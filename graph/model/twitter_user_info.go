package model

type TwitterUserInfo struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	DisplayName     string `json:"displayName"`
	ProfileImageURL string `json:"profileImageURL"`
	Verified        bool   `json:"verified"`
	Invited         bool   `json:"invited"`
}
