package remote

import "fmt"

// define the remote session
type RemoteSession struct {
    token string
    tokenVerified bool
}

func (s *RemoteSession) GetToken() string {
    return s.token
}

// create a new remote session and return it
func NewRemoteSession(token string) (*RemoteSession, error) {
    if token == "" {
        return nil, fmt.Errorf("token not provided")
    }

    session := &RemoteSession{
        token: token,
        tokenVerified: false,
    }

    // authenticate the session
    if err := session.Authenticate(); err != nil {
        return nil, err
    }

    return session, nil
}
