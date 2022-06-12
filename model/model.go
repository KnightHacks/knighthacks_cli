package model

type LoginPayload struct {
	// If false then you must register immediately following this. Else, you are logged in and have access to your own user.
	AccountExists bool    `json:"accountExists"`
	User          *User   `json:"user"`
	Jwt           *string `json:"jwt"`
}

type NewUser struct {
	FirstName   string         `json:"firstName"`
	LastName    string         `json:"lastName"`
	Email       string         `json:"email"`
	PhoneNumber string         `json:"phoneNumber"`
	Pronouns    *PronounsInput `json:"pronouns"`
	Age         *int           `json:"age"`
}

type OAuth struct {
	Provider Provider `json:"provider"`
	UID      string   `json:"uid"`
}

// Example:
// subjective=he
// objective=him
type Pronouns struct {
	Subjective string `json:"subjective"`
	Objective  string `json:"objective"`
}

type PronounsInput struct {
	Subjective string `json:"subjective"`
	Objective  string `json:"objective"`
}

type User struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	FullName    string    `json:"fullName"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	Pronouns    *Pronouns `json:"pronouns"`
	Age         *int      `json:"age"`
	OAuth       *OAuth    `json:"oAuth"`
}

type Provider string

const (
	ProviderGithub Provider = "GITHUB"
	ProviderGmail  Provider = "GMAIL"
)
