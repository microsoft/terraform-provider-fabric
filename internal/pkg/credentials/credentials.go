package credentials

import (
	"errors"

	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

// Credentials is the base structure for all credential types
type Credentials struct {
	Type           fabcore.CredentialType
	CredentialData []CredentialEntry `json:"credentialData"`
}

type CredentialEntry struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// NewAnonymousCredentials creates a new instance of AnonymousCredentials
func NewAnonymousCredentials() (*Credentials, error) {
	creds := Credentials{
		Type:           fabcore.CredentialTypeAnonymous,
		CredentialData: make([]CredentialEntry, 0),
	}

	return &creds, nil
}

// NewOAuth2Credentials creates a new instance of OAuth2Credentials
func NewOAuth2Credentials(accessToken string) (*Credentials, error) {
	if accessToken == "" {
		return nil, errors.New("accessToken is required")
	}

	creds := Credentials{
		Type:           fabcore.CredentialTypeOAuth2,
		CredentialData: make([]CredentialEntry, 0),
	}

	creds.CredentialData = append(creds.CredentialData, CredentialEntry{
		Name:  "accessToken",
		Value: accessToken,
	})

	return &creds, nil
}

// NewWindowsCredentials creates a new instance of WindowsCredentials
func NewWindowsCredentials(username, password string) (*Credentials, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	creds := Credentials{
		Type:           fabcore.CredentialTypeWindows,
		CredentialData: make([]CredentialEntry, 0),
	}

	creds.CredentialData = append(creds.CredentialData, CredentialEntry{
		Name:  "username",
		Value: username,
	})

	creds.CredentialData = append(creds.CredentialData, CredentialEntry{
		Name:  "password",
		Value: password,
	})

	return &creds, nil
}

// NewKeyCredentials creates a new instance of KeyCredentials
func NewKeyCredentials(key string) (*Credentials, error) {
	if key == "" {
		return nil, errors.New("key is required")
	}

	creds := Credentials{
		Type:           fabcore.CredentialTypeKey,
		CredentialData: make([]CredentialEntry, 0),
	}

	creds.CredentialData = append(creds.CredentialData, CredentialEntry{
		Name:  "key",
		Value: key,
	})

	return &creds, nil
}

// NewBasicCredentials creates a new instance of BasicCredentials
func NewBasicCredentials(username, password string) (*Credentials, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	creds := Credentials{
		Type:           fabcore.CredentialTypeBasic,
		CredentialData: make([]CredentialEntry, 0),
	}

	creds.CredentialData = append(creds.CredentialData, CredentialEntry{
		Name:  "username",
		Value: username,
	})

	creds.CredentialData = append(creds.CredentialData, CredentialEntry{
		Name:  "password",
		Value: password,
	})

	return &creds, nil
}

func NewCredentials(creds fabcore.CredentialsClassification) (*Credentials, error) {
	if creds == nil {
		return nil, errors.New("creds is required")
	}

	credentialType := creds.GetCredentials()

	switch *credentialType.CredentialType {
	case fabcore.CredentialTypeAnonymous:
		return NewAnonymousCredentials()
	case fabcore.CredentialTypeOAuth2:
		return NewOAuth2Credentials("TODO")
	case fabcore.CredentialTypeWindows:
		c := creds.(*fabcore.WindowsCredentials)

		return NewWindowsCredentials(*c.Username, *c.Password)

	case fabcore.CredentialTypeKey:
		c := creds.(*fabcore.KeyCredentials)

		return NewKeyCredentials(*c.Key)

	case fabcore.CredentialTypeBasic:
		c := creds.(*fabcore.BasicCredentials)

		return NewBasicCredentials(*c.Username, *c.Password)

	default:
		return nil, errors.New("unsupported credential type")
	}
}
