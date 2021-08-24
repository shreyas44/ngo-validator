package main

type VerificationStatus string

const (
	StatusSuccess VerificationStatus = "success"
	StatusFailure VerificationStatus = "fail"
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
