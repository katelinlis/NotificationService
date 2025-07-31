package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/avenir/notification-service/internal/domain/model"
)

type UserClient struct {
	baseURL string
	client  *http.Client
}

func NewUserClient(baseURL string, client *http.Client) *UserClient {
	return &UserClient{
		baseURL: baseURL,
		client:  client,
	}
}
func (uc *UserClient) GetUser(ctx context.Context, userID int) (*model.User, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/users/%d", uc.baseURL, userID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := uc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user: %s", resp.Status)
	}

	var user model.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
func (uc *UserClient) GetUsers(ctx context.Context, userIDs []int) ([]model.User, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/users/batch", uc.baseURL), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for _, id := range userIDs {
		q.Add("id", fmt.Sprintf("%d", id))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := uc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get users: %s", resp.Status)
	}

	var users []model.User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}
