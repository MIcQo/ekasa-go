package ekasa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Client is the eKasa API client.
type Client struct {
	publicKey  string
	privateKey string
	opts       clientOptions
	signer     *signer
}

// NewClient creates a new eKasa API client.
func NewClient(publicKey, privateKey string, opts ...ClientOption) *Client {
	co := defaultOptions()
	for _, opt := range opts {
		opt(&co)
	}

	return &Client{
		publicKey:  publicKey,
		privateKey: privateKey,
		opts:       co,
		signer:     newSigner(publicKey, privateKey),
	}
}

// Registrations

// RegisterReceipt submits a request to register a receipt.
func (c *Client) RegisterReceipt(ctx context.Context, receipt *CreateReceiptRegistrationDto) (*ReceiptRegistrationDto, error) {
	var result ReceiptRegistrationDto
	err := c.doRequest(ctx, http.MethodPost, "/v1/registrations/receipts", receipt, &result, true)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelReceipt cancels a receipt registration request.
func (c *Client) CancelReceipt(ctx context.Context, cashRegisterCode, externalId string) (*ReceiptRegistrationStateChangeResultDto, error) {
	path := fmt.Sprintf("/v1/registrations/receipts/cancel?cashRegisterCode=%s&externalId=%s", url.QueryEscape(cashRegisterCode), url.QueryEscape(externalId))
	var result ReceiptRegistrationStateChangeResultDto
	err := c.doRequest(ctx, http.MethodPost, path, nil, &result, true)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Customers

// GetCustomers retrieves a list of customers based on the provided filter.
func (c *Client) GetCustomers(ctx context.Context, filter *CustomerFilterDto) (*GetCustomerListResultDto, error) {
	path := "/v1/customers"
	if filter != nil {
		params := url.Values{}
		if len(filter.IDs) > 0 {
			for _, id := range filter.IDs {
				params.Add("ids", id)
			}
		}
		if filter.ExternalID != "" {
			params.Set("externalId", filter.ExternalID)
		}
		if filter.ModifiedAfter != nil {
			params.Set("modifiedAfter", filter.ModifiedAfter.Format("2006-01-02T15:04:05.000000Z"))
		}
		if filter.IsActive != nil {
			if *filter.IsActive {
				params.Set("isActive", "true")
			} else {
				params.Set("isActive", "false")
			}
		}
		if filter.CardID != "" {
			params.Set("cardId", filter.CardID)
		}
		if len(filter.CardSerialNumbers) > 0 {
			for _, sn := range filter.CardSerialNumbers {
				params.Add("cardSerialNumbers", sn)
			}
		}
		if len(params) > 0 {
			path += "?" + params.Encode()
		}
	}

	var result GetCustomerListResultDto
	err := c.doRequest(ctx, http.MethodGet, path, nil, &result, true)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCustomer retrieves a specific customer by ID.
func (c *Client) GetCustomer(ctx context.Context, id string) (*CustomerDto, error) {
	path := "/v1/customers/" + url.PathEscape(id)
	var result CustomerDto
	err := c.doRequest(ctx, http.MethodGet, path, nil, &result, true)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// FindCustomer attempts to find a customer by ID, returning nil if not found (404).
func (c *Client) FindCustomer(ctx context.Context, id string) (*CustomerDto, error) {
	customer, err := c.GetCustomer(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, nil
		}
		return nil, err
	}
	return customer, nil
}

// Internal

func (c *Client) doRequest(ctx context.Context, method, path string, payload interface{}, result interface{}, sign bool) error {
	fullURL := strings.TrimSuffix(c.opts.baseURL, "/") + "/" + strings.TrimPrefix(path, "/")

	var body []byte
	var err error
	if payload != nil {
		body, err = json.Marshal(payload)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bytes.NewReader(body))
	if err != nil {
		return err
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json; charset=utf-8"
	headers["Accept"] = "application/json"
	if c.opts.tenantID != "" {
		headers["__tenant"] = c.opts.tenantID
	}

	if sign {
		if err := c.signer.sign(method, fullURL, headers, body); err != nil {
			return err
		}
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.opts.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return c.handleError(resp)
	}

	if result != nil {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return json.Unmarshal(respBody, result)
	}

	return nil
}

func (c *Client) handleError(resp *http.Response) error {
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 400 || resp.StatusCode == 422 {
		var vpd ValidationProblemDetails
		if err := json.Unmarshal(respBody, &vpd); err == nil {
			return fmt.Errorf("validation error (%d): %s, errors: %+v", resp.StatusCode, vpd.Title, vpd.Errors)
		}
	}

	var pd ProblemDetails
	if err := json.Unmarshal(respBody, &pd); err == nil {
		return fmt.Errorf("api error (%d): %s, detail: %s", resp.StatusCode, pd.Title, pd.Detail)
	}

	return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
}
