package main

type VerificationStatus string

const (
	StatusFail        VerificationStatus = "fail"
	StatusSuccess     VerificationStatus = "success"
	StatusUnfulfilled VerificationStatus = "unfulfilled"
)

type NGOInfo struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	State   string `json:"state"`
	PAN     string `json:"pan"`
}

type VerifyInput NGOInfo

type VerifyResponse struct {
	VerificationStatus VerificationStatus `json:"verification_status"`
	Message            string             `json:"message"`
}
