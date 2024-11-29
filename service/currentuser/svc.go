package currentuser

import (
	"context"
	"time"

	"github.com/Gardego5/garrettdavis.dev/model"
	"github.com/Gardego5/garrettdavis.dev/utils/bimarshal"
	"github.com/google/go-github/v66/github"
)

type Service struct {
	users        bimarshal.Cache[github.User]
	accessTokens bimarshal.Cache[model.GHAccessToken]
}

func New(caches bimarshal.RegisteredCaches) *Service {
	return &Service{
		users:        bimarshal.Get[github.User](caches),
		accessTokens: bimarshal.Get[model.GHAccessToken](caches),
	}
}

func (s *Service) GetUserByAccessToken(
	ctx context.Context,
	accessToken string,
) (*github.User, error) {
	return s.users.GetOrSet(ctx, accessToken, func() (*github.User, time.Duration, error) {
		user, _, err := github.NewClient(nil).
			WithAuthToken(accessToken).Users.Get(ctx, "")
		return user, 1 * time.Hour, err
	})
}

func (s *Service) GetUserBySession(
	ctx context.Context,
	session string,
) (*github.User, error) {
	oauth, err := s.accessTokens.Get(ctx, session)
	if err != nil {
		return nil, err
	}

	return s.GetUserByAccessToken(ctx, oauth.AccessToken)
}
