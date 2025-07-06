package response

type GenerateInvoicePDFResponse struct {
	FilePath string `json:"file_path"`
	// Optionally, add a download URL if serving via HTTP
	// DownloadURL string `json:"download_url,omitempty"`
}
