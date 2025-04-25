package dto

type AccountResponse struct {
	NoRekening string `json:"no_rekening"`
}

type BalanceResponse struct {
	Saldo float64 `json:"saldo"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Remark string `json:"remark"`
}

type AccountRegistration struct {
	Nama string `json:"nama"`
	NIK  string `json:"nik"`
	NoHP string `json:"no_hp"`
}

type WithdrawDepositRequest struct {
	NoRekening string  `json:"no_rekening"`
	Amount     float64 `json:"saldo"`
}
