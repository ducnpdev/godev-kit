package request

type InvoiceItem struct {
	Description string `json:"description"`
	UnitCost    string `json:"unit_cost"`
	Qty         string `json:"qty"`
	Amount      string `json:"amount"`
}

type GenerateInvoicePDFRequest struct {
	Number      string        `json:"number"`
	Date        string        `json:"date"`
	BilledTo    []string      `json:"billed_to"`
	CompanyInfo []string      `json:"company_info"`
	Items       []InvoiceItem `json:"items"`
	Subtotal    string        `json:"subtotal"`
	Discount    string        `json:"discount"`
	TaxRate     string        `json:"tax_rate"`
	Tax         string        `json:"tax"`
	Total       string        `json:"total"`
	Terms       string        `json:"terms"`
	BankDetails []string      `json:"bank_details"`
}
