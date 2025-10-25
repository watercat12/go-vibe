package bank_link

import (
	"bytes"
	"context"
	"encoding/json"
	"e-wallet/internal/domain/bank_link"
	"e-wallet/internal/ports"
	"net/http"
)

type bankLinkClient struct {
	baseURL string
}

func NewBankLinkClient(baseURL string) ports.BankLinkClient {
	return &bankLinkClient{baseURL: baseURL}
}

type linkAccountRequest struct {
	BankCode    string `json:"bank_code"`
	AccountType string `json:"account_type"`
}

type linkAccountResponse struct {
	AccessToken  string `json:"access_token"`
	BankCode     string `json:"bank_code"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (c *bankLinkClient) LinkAccount(ctx context.Context, bankCode, accountType string) (*bank_link.BankLink, error) {
	reqBody := linkAccountRequest{
		BankCode:    bankCode,
		AccountType: accountType,
	}

	jsonData, _ := json.Marshal(reqBody)

	req, _ := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/mock/bbecc26a-c96a-4aa6-8178-0ee9cde9f390/303/bank-link", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &BankLinkError{Message: "failed to link account"}
	}

	var res linkAccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return bank_link.NewBankLink("", bankCode, accountType, res.AccessToken, res.RefreshToken, res.ExpiresIn), nil
}