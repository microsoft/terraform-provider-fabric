package credentials

import (
	"errors"

	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

// Credentials is the base structure for all credential types
type Credentials struct {
	credentialType fabcore.CredentialType
	CredentialData []CredentialEntry `json:"credentialData"`
}

type CredentialEntry struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// NewAnonymousCredentials creates a new instance of AnonymousCredentials
func NewAnonymousCredentials() (*Credentials, error) {
	creds := Credentials{
		credentialType: fabcore.CredentialTypeAnonymous,
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
		credentialType: fabcore.CredentialTypeOAuth2,
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
		credentialType: fabcore.CredentialTypeWindows,
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
		credentialType: fabcore.CredentialTypeKey,
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
		credentialType: fabcore.CredentialTypeBasic,
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
