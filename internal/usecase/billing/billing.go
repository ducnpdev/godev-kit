package billing

import (
	"github.com/signintech/gopdf"
)

type InvoiceItem struct {
	Description string
	UnitCost    string
	Qty         string
	Amount      string
}

type InvoiceData struct {
	Number      string
	Date        string
	BilledTo    []string
	CompanyInfo []string
	Items       []InvoiceItem
	Subtotal    string
	Discount    string
	TaxRate     string
	Tax         string
	Total       string
	Terms       string
	BankDetails []string
}

const (
	marginLeft     = 50.0
	marginTop      = 50.0
	lineHeight     = 20.0
	pageWidth      = 595.28 // A4 width in points
	pageHeight     = 841.89 // A4 height in points
	tableRowHeight = 18.0
)

type UseCase struct{}

func New() *UseCase {
	return &UseCase{}
}

// GenerateInvoicePDF generates a PDF invoice and returns the file path
func (uc *UseCase) GenerateInvoicePDF(data InvoiceData, outputPath string) error {
	pdf := gopdf.GoPdf{}
	mm6ToPx := 22.68

	pdf.Start(gopdf.Config{
		PageSize: *gopdf.PageSizeA4,
		TrimBox:  gopdf.Box{Left: mm6ToPx, Top: mm6ToPx, Right: gopdf.PageSizeA4.W - mm6ToPx, Bottom: gopdf.PageSizeA4.H - mm6ToPx},
	})
	opt := gopdf.PageOption{
		PageSize: gopdf.PageSizeA4,
		TrimBox:  &gopdf.Box{Left: mm6ToPx, Top: mm6ToPx, Right: gopdf.PageSizeA4.W - mm6ToPx, Bottom: gopdf.PageSizeA4.H - mm6ToPx},
	}
	pdf.AddPageWithOption(opt)

	if err := pdf.AddTTFFont("roboto", "./docs/front/Roboto-Regular.ttf"); err != nil {
		return err
	}
	if err := pdf.AddTTFFont("roboto-bold", "./docs/front/Roboto-Regular.ttf"); err != nil {
		return err
	}

	headerBottomY := drawHeader(&pdf, data)
	tableBottomY := drawTable(&pdf, data.Items, headerBottomY)
	summaryBottomY := drawSummary(&pdf, data, tableBottomY)
	drawFooter(&pdf, data, summaryBottomY)

	if err := pdf.WritePdf(outputPath); err != nil {
		return err
	}
	return nil
}

// --- PDF Drawing Functions (copied from main.go, made unexported) ---

func drawHeader(pdf *gopdf.GoPdf, data InvoiceData) float64 {
	// Title
	pdf.SetFont("roboto-bold", "", 28) // For bold
	pdf.SetTextColor(30, 60, 120)
	pdf.SetX(marginLeft)
	pdf.SetY(marginTop)
	pdf.Cell(nil, "Invoice")

	// Logo (placeholder rectangle)
	logoSize := 80.0
	logoX := pageWidth - marginLeft - logoSize
	logoY := marginTop
	pdf.SetLineWidth(0.5)
	pdf.RectFromUpperLeftWithStyle(logoX, logoY, logoSize, logoSize, "D")
	pdf.SetFont("roboto", "", 11)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetX(logoX + 10)
	pdf.SetY(logoY + logoSize/2)
	pdf.Cell(nil, "YOUR LOGO")

	// Invoice number and date of issue (2 columns)
	topInfoY := marginTop + 45
	col1X := marginLeft
	col2X := marginLeft + 180

	pdf.SetFont("roboto-bold", "", 11)
	pdf.SetTextColor(30, 60, 120)
	pdf.SetX(col1X)
	pdf.SetY(topInfoY)
	pdf.Cell(nil, "INVOICE NUMBER:")
	pdf.SetX(col2X)
	pdf.Cell(nil, "DATE OF ISSUE:")

	pdf.SetFont("roboto", "", 11)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetX(col1X)
	pdf.SetY(topInfoY + 15)
	pdf.Cell(nil, data.Number)
	pdf.SetX(col2X)
	pdf.Cell(nil, data.Date)

	// Billed to, Company name, Company contact (3 columns)
	sectionY := topInfoY + 40
	col3X := marginLeft + 350

	pdf.SetFont("roboto-bold", "", 11)
	pdf.SetTextColor(30, 60, 120)
	pdf.SetX(col1X)
	pdf.SetY(sectionY)
	pdf.Cell(nil, "BILLED TO")
	pdf.SetX(col2X)
	pdf.Cell(nil, "YOUR COMPANY NAME")
	pdf.SetX(col3X)
	pdf.Cell(nil, "")

	// Details under each column
	pdf.SetFont("roboto", "", 10)
	pdf.SetTextColor(0, 0, 0)
	billedTo := []string{"Client name", "123 Your Street", "City,State, Country", "Zip Code", "Phone"}
	companyInfo := []string{"Building name", "123 Your Street", "City,State, Country", "Zip Code", "Phone"}
	companyContact := []string{"+1-541-754-3010", "you@email.com", "yourwebsite.com"}

	for i := 0; i < 5; i++ {
		pdf.SetX(col1X)
		pdf.SetY(sectionY + 15 + float64(i)*13)
		pdf.Cell(nil, billedTo[i])
		pdf.SetX(col2X)
		pdf.Cell(nil, companyInfo[i])
		if i < len(companyContact) {
			pdf.SetX(col3X)
			pdf.Cell(nil, companyContact[i])
		}
	}
	// Return the Y position after the last line
	return sectionY + 15 + float64(5)*13 + 10 // +10 for extra spacing
}

func drawTable(pdf *gopdf.GoPdf, items []InvoiceItem, startY float64) float64 {
	tableTop := startY
	tableLeft := marginLeft
	tableWidth := pageWidth - 2*marginLeft
	tableRowHeight := 18.0
	tableColWidths := []float64{200, 100, 100, 100}

	// Table header background
	pdf.SetFillColor(240, 245, 250)
	pdf.RectFromUpperLeftWithStyle(tableLeft, tableTop, tableWidth, tableRowHeight, "F")

	// Table headers
	pdf.SetFont("roboto-bold", "", 11)
	pdf.SetTextColor(30, 60, 120)
	pdf.SetY(tableTop + 4)
	pdf.SetX(tableLeft + 8)
	pdf.Cell(nil, "Description")
	pdf.SetX(tableLeft + tableColWidths[0] + 8)
	pdf.Cell(nil, "Unit cost")
	pdf.SetX(tableLeft + tableColWidths[0] + tableColWidths[1] + 8)
	pdf.Cell(nil, "QTY/HR Rate")
	pdf.SetX(tableLeft + tableColWidths[0] + tableColWidths[1] + tableColWidths[2] + 8)
	pdf.Cell(nil, "Amount")

	// Table header bottom border
	pdf.SetStrokeColor(200, 200, 200)
	pdf.Line(tableLeft, tableTop+tableRowHeight, tableLeft+tableWidth, tableTop+tableRowHeight)

	// Table rows
	pdf.SetFont("roboto", "", 10)
	pdf.SetTextColor(0, 0, 0)
	rowY := tableTop + tableRowHeight
	for _, item := range items {
		pdf.SetY(rowY + 4)
		pdf.SetX(tableLeft + 8)
		pdf.Cell(nil, item.Description)
		pdf.SetX(tableLeft + tableColWidths[0] + 8)
		pdf.Cell(nil, item.UnitCost)
		pdf.SetX(tableLeft + tableColWidths[0] + tableColWidths[1] + 8)
		pdf.Cell(nil, item.Qty)
		pdf.SetX(tableLeft + tableColWidths[0] + tableColWidths[1] + tableColWidths[2] + 8)
		pdf.Cell(nil, item.Amount)
		rowY += tableRowHeight
	}

	// Table bottom border
	pdf.SetStrokeColor(200, 200, 200)
	pdf.Line(tableLeft, rowY, tableLeft+tableWidth, rowY)

	return rowY + 10 // +10 for extra spacing
}

func drawSummary(pdf *gopdf.GoPdf, data InvoiceData, startY float64) float64 {
	summaryLeft := pageWidth - marginLeft - 200
	summaryY := startY

	pdf.SetFont("roboto-bold", "", 11)
	pdf.SetTextColor(30, 60, 120)
	pdf.SetX(summaryLeft)
	pdf.SetY(summaryY)
	pdf.Cell(nil, "Subtotal:")
	pdf.SetX(summaryLeft + 120)
	pdf.Cell(nil, data.Subtotal)

	pdf.SetX(summaryLeft)
	pdf.SetY(summaryY + 18)
	pdf.Cell(nil, "Discount:")
	pdf.SetX(summaryLeft + 120)
	pdf.Cell(nil, data.Discount)

	pdf.SetX(summaryLeft)
	pdf.SetY(summaryY + 36)
	pdf.Cell(nil, "Tax Rate:")
	pdf.SetX(summaryLeft + 120)
	pdf.Cell(nil, data.TaxRate)

	pdf.SetX(summaryLeft)
	pdf.SetY(summaryY + 54)
	pdf.Cell(nil, "Tax:")
	pdf.SetX(summaryLeft + 120)
	pdf.Cell(nil, data.Tax)

	pdf.SetFont("roboto-bold", "", 13)
	pdf.SetTextColor(30, 60, 120)
	pdf.SetX(summaryLeft)
	pdf.SetY(summaryY + 80)
	pdf.Cell(nil, "Total:")
	pdf.SetX(summaryLeft + 120)
	pdf.Cell(nil, data.Total)

	return summaryY + 110 // +30 for extra spacing
}

func drawFooter(pdf *gopdf.GoPdf, data InvoiceData, startY float64) {
	footerY := startY
	pdf.SetFont("roboto", "", 10)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetX(marginLeft)
	pdf.SetY(footerY)
	pdf.Cell(nil, data.Terms)

	// Bank details (if any)
	if len(data.BankDetails) > 0 {
		pdf.SetY(footerY + 30)
		pdf.SetFont("roboto-bold", "", 11)
		pdf.SetTextColor(30, 60, 120)
		pdf.SetX(marginLeft)
		pdf.Cell(nil, "Bank Details:")
		pdf.SetFont("roboto", "", 10)
		pdf.SetTextColor(0, 0, 0)
		for i, line := range data.BankDetails {
			pdf.SetY(footerY + 50 + float64(i)*13)
			pdf.SetX(marginLeft)
			pdf.Cell(nil, line)
		}
	}
}
