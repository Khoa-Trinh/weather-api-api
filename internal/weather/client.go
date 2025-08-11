package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	apiKey     string
	unitGroup  string
	baseURL    string
	httpClient *http.Client
}

func NewClient(apiKey, defaultUnitGroup string, timeout time.Duration) *Client {
	return &Client{
		apiKey:     apiKey,
		unitGroup:  defaultUnitGroup,
		baseURL:    "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline",
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (c *Client) Fetch(ctx context.Context, place, unitGroup string) ([]byte, int, error) {
	if unitGroup == "" {
		unitGroup = c.unitGroup
	}
	u := fmt.Sprintf("%s/%s", c.baseURL, url.PathEscape(place))
	q := url.Values{}
	q.Set("unitGroup", unitGroup)
	q.Set("key", c.apiKey)
	q.Set("contentType", "json")
	full := u + "?" + q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, full, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", "weather-api/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		msg := map[string]any{
			"error":       "upstream_error",
			"status_code": resp.StatusCode,
			"provider":    "visual_crossing",
		}
		_ = json.Unmarshal(body, &msg)
		out, _ := json.Marshal(msg)
		return out, resp.StatusCode, fmt.Errorf("upstream status %d", resp.StatusCode)
	}
	return body, resp.StatusCode, nil
}
