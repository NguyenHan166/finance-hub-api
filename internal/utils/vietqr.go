package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// VietQRBank represents a bank from VietQR API
type VietQRBank struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Code              string `json:"code"`
	BIN               string `json:"bin"`
	ShortName         string `json:"shortName"`
	Logo              string `json:"logo"`
	TransferSupported int    `json:"transferSupported"`
	LookupSupported   int    `json:"lookupSupported"`
}

// VietQRResponse represents the API response
type VietQRResponse struct {
	Code string         `json:"code"`
	Desc string         `json:"desc"`
	Data []VietQRBank   `json:"data"`
}

// VietQRService handles VietQR API integration
type VietQRService struct {
	apiURL string
	client *http.Client
}

// NewVietQRService creates a new VietQR service
func NewVietQRService() *VietQRService {
	return &VietQRService{
		apiURL: "https://api.vietqr.io/v2",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetBanks retrieves list of all banks
func (s *VietQRService) GetBanks() ([]VietQRBank, error) {
	url := fmt.Sprintf("%s/banks", s.apiURL)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch banks: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("VietQR API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var vietQRResp VietQRResponse
	if err := json.Unmarshal(body, &vietQRResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if vietQRResp.Code != "00" {
		return nil, fmt.Errorf("VietQR API error: %s", vietQRResp.Desc)
	}

	return vietQRResp.Data, nil
}

// GetBankByCode retrieves a specific bank by code
func (s *VietQRService) GetBankByCode(code string) (*VietQRBank, error) {
	banks, err := s.GetBanks()
	if err != nil {
		return nil, err
	}

	for _, bank := range banks {
		if bank.Code == code {
			return &bank, nil
		}
	}

	return nil, fmt.Errorf("bank with code %s not found", code)
}

// GetBankByBIN retrieves a specific bank by BIN
func (s *VietQRService) GetBankByBIN(bin string) (*VietQRBank, error) {
	banks, err := s.GetBanks()
	if err != nil {
		return nil, err
	}

	for _, bank := range banks {
		if bank.BIN == bin {
			return &bank, nil
		}
	}

	return nil, fmt.Errorf("bank with BIN %s not found", bin)
}

// SearchBanks searches banks by name or code
func (s *VietQRService) SearchBanks(query string) ([]VietQRBank, error) {
	banks, err := s.GetBanks()
	if err != nil {
		return nil, err
	}

	var results []VietQRBank
	queryLower := toLowerASCII(query)

	for _, bank := range banks {
		if contains(toLowerASCII(bank.Name), queryLower) ||
			contains(toLowerASCII(bank.Code), queryLower) ||
			contains(toLowerASCII(bank.ShortName), queryLower) {
			results = append(results, bank)
		}
	}

	return results, nil
}

// Helper functions
func toLowerASCII(s string) string {
	result := make([]rune, len([]rune(s)))
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			result[i] = r + 32
		} else {
			result[i] = r
		}
	}
	return string(result)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || search(s, substr))
}

func search(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
