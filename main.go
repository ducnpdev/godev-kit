package main

import (
	"log"

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

func main() {
	// Prepare invoice data
	data := InvoiceData{
		Number: "00001",
		Date:   "MM/DD/YYYY",
		Items: []InvoiceItem{
			{"Your item name", "$0.00", "1", "$0.00"},
			{"Your item name", "$0.00", "1", "$0.00"},
			{"Your item name", "$0.00", "1", "$0.00"},
			{"Your item name", "$0.00", "1", "$0.00"},
			{"Your item name", "$0.00", "1", "$0.00"},
			{"Your item name", "$0.00", "1", "$0.00"},
			{"Your item name", "$0.00", "1", "$0.00"},
			{"Your item name", "$0.00", "1", "$0.00"},
		},
		Subtotal: "$0.00",
		Discount: "$0.00",
		TaxRate:  "0 %",
		Tax:      "$0.00",
		Total:    "$0.00",
		Terms:    "Please pay invoice by MM/DD/YYYY",
	}

	// Setup PDF
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
		log.Print(err.Error())
		return
	}

	// Draw sections
	drawHeader(&pdf, data)
	drawTable(&pdf, data.Items)
	drawSummary(&pdf, data)
	drawFooter(&pdf, data)

	pdf.WritePdf("hello.pdf")
}

func drawHeader(pdf *gopdf.GoPdf, data InvoiceData) {
	pdf.SetFont("roboto", "", 24)
	pdf.SetX(marginLeft)
	pdf.SetY(marginTop)
	pdf.Cell(nil, "Invoice")

	// Invoice number and date
	pdf.SetFont("roboto", "", 10)
	pdf.SetX(marginLeft)
	pdf.SetY(marginTop + lineHeight)
	pdf.Cell(nil, "INVOICE NUMBER:")
	pdf.SetX(marginLeft + 110)
	pdf.Cell(nil, data.Number)
	pdf.SetX(marginLeft + 200)
	pdf.Cell(nil, "DATE OF ISSUE:")
	pdf.SetX(marginLeft + 300)
	pdf.Cell(nil, data.Date)

	// Logo (placeholder rectangle)
	pdf.SetLineWidth(0.5)
	logoX := pageWidth - marginLeft - 80
	logoY := marginTop
	pdf.RectFromUpperLeftWithStyle(logoX, logoY, 70, 70, "D")
	pdf.SetX(logoX + 10)
	pdf.SetY(logoY + 30)
	pdf.SetFont("roboto", "", 10)
	pdf.Cell(nil, "YOUR LOGO")

	// Billed to and Company info
	pdf.SetFont("roboto", "", 10)
	pdf.SetX(marginLeft)
	pdf.SetY(marginTop + 2*lineHeight)
	pdf.Cell(nil, "BILLED TO")
	pdf.SetX(marginLeft + 180)
	pdf.Cell(nil, "YOUR COMPANY NAME")

	// Example details
	billedTo := []string{"Client name", "123 Your Street", "City,State, Country", "Zip Code", "Phone"}
	companyInfo := []string{"Building name", "123 Your Street", "City,State, Country", "Zip Code", "Phone"}
	companyContact := []string{"+1-541-754-3010", "you@email.com", "yourwebsite.com"}

	for i := 0; i < len(billedTo); i++ {
		pdf.SetX(marginLeft)
		pdf.SetY(marginTop + 2*lineHeight + 15 + float64(i)*15)
		pdf.Cell(nil, billedTo[i])
		pdf.SetX(marginLeft + 180)
		pdf.Cell(nil, companyInfo[i])
		if i < len(companyContact) {
			pdf.SetX(marginLeft + 350)
			pdf.Cell(nil, companyContact[i])
		}
	}
}

func drawTable(pdf *gopdf.GoPdf, items []InvoiceItem) {
	tableTop := marginTop + 2*lineHeight + 100
	tableLeft := marginLeft
	tableWidth := pageWidth - 2*marginLeft
	tableRowHeight := 18.0
	tableColWidths := []float64{200, 100, 100, 100}

	// Table header background
	pdf.SetFillColor(240, 245, 250)
	pdf.RectFromUpperLeftWithStyle(tableLeft, tableTop, tableWidth, tableRowHeight, "F")

	// Table headers
	pdf.SetFont("roboto", "", 12)
	pdf.SetTextColor(80, 90, 100)
	pdf.SetX(tableLeft + 5)
	pdf.SetY(tableTop + 4)
	pdf.Cell(nil, "Description")
	pdf.SetX(tableLeft + tableColWidths[0] + 5)
	pdf.Cell(nil, "Unit cost")
	pdf.SetX(tableLeft + tableColWidths[0] + tableColWidths[1] + 5)
	pdf.Cell(nil, "QTY/HR Rate")
	pdf.SetX(tableLeft + tableColWidths[0] + tableColWidths[1] + tableColWidths[2] + 5)
	pdf.Cell(nil, "Amount")

	// Table rows
	pdf.SetFont("roboto", "", 11)
	pdf.SetTextColor(0, 0, 0)
	rowY := tableTop + tableRowHeight
	for _, item := range items {
		pdf.SetX(tableLeft + 5)
		pdf.SetY(rowY + 4)
		pdf.Cell(nil, item.Description)
		pdf.SetX(tableLeft + tableColWidths[0] + 5)
		pdf.Cell(nil, item.UnitCost)
		pdf.SetX(tableLeft + tableColWidths[0] + tableColWidths[1] + 5)
		pdf.Cell(nil, item.Qty)
		pdf.SetX(tableLeft + tableColWidths[0] + tableColWidths[1] + tableColWidths[2] + 5)
		pdf.Cell(nil, item.Amount)
		rowY += tableRowHeight
	}

	// Table border lines
	pdf.SetLineWidth(0.5)
	borderY := tableTop
	for i := 0; i <= len(items)+1; i++ {
		pdf.Line(tableLeft, borderY, tableLeft+tableWidth, borderY)
		borderY += tableRowHeight
	}
	// Vertical lines
	colX := tableLeft
	for _, w := range tableColWidths {
		pdf.Line(colX, tableTop, colX, tableTop+tableRowHeight*float64(len(items)+1))
		colX += w
	}
	pdf.Line(colX, tableTop, colX, tableTop+tableRowHeight*float64(len(items)+1))
}

func drawSummary(pdf *gopdf.GoPdf, data InvoiceData) {
	tableTop := marginTop + 2*lineHeight + 100
	tableRowHeight := 18.0
	tableColWidths := []float64{200, 100, 100, 100}
	summaryY := tableTop + tableRowHeight*float64(len(data.Items)+2)
	tableLeft := marginLeft

	pdf.SetFont("roboto", "", 11)
	pdf.SetX(tableLeft + tableColWidths[0] + tableColWidths[1] + tableColWidths[2] - 10)
	pdf.SetY(summaryY)
	pdf.Cell(nil, "Subtotal")
	pdf.SetX(tableLeft + tableColWidths[0] + tableColWidths[1] + tableColWidths[2] + 30)
	pdf.Cell(nil, data.Subtotal)

	pdf.SetX(tableLeft + tableColWidths[0] + tableColWidths[1] + tableColWidths[2] - 10)
	pdf.SetY(summaryY + 15)
	pdf.Cell(nil, "Discount")
	pdf.SetX(tableLeft + tableColWidths[0] + tableColWidths[1] + tableColWidths[2] + 30)
	pdf.Cell(nil, data.Discount)

	pdf.SetX(tableLeft + tableColWidths[0] + tableColWidths[1] + tableColWidths[2] - 10)
	pdf.SetY(summaryY + 30)
	pdf.Cell(nil, "Tax rate")
	pdf.SetX(tableLeft + tableColWidths[0] + tableColWidths[1] + tableColWidths[2] + 30)
	pdf.Cell(nil, data.TaxRate)

	pdf.SetX(tableLeft + tableColWidths[0] + tableColWidths[1] + tableColWidths[2] - 10)
	pdf.SetY(summaryY + 45)
	pdf.Cell(nil, "Tax")
	pdf.SetX(tableLeft + tableColWidths[0] + tableColWidths[1] + tableColWidths[2] + 30)
	pdf.Cell(nil, data.Tax)
}

func drawFooter(pdf *gopdf.GoPdf, data InvoiceData) {
	tableTop := marginTop + 2*lineHeight + 100
	tableRowHeight := 18.0
	footerY := tableTop + tableRowHeight*float64(len(data.Items)+2) + 80

	pdf.SetFont("roboto", "", 10)
	pdf.SetX(marginLeft)
	pdf.SetY(footerY)
	pdf.Cell(nil, "TERMS")
	pdf.SetX(marginLeft + 150)
	pdf.Cell(nil, "BANK ACCOUNT DETAILS")
	pdf.SetX(pageWidth - marginLeft - 100)
	pdf.Cell(nil, "INVOICE TOTAL")

	pdf.SetFont("roboto", "", 9)
	pdf.SetX(marginLeft)
	pdf.SetY(footerY + 15)
	pdf.Cell(nil, data.Terms)

	// Example bank details
	bankDetails := []string{
		"Account Holder:",
		"Account number:",
		"ABA rtn: 026073150",
		"Wire rtn: 026073008",
	}
	for i, line := range bankDetails {
		pdf.SetX(marginLeft + 150)
		pdf.SetY(footerY + 15 + float64(i)*15)
		pdf.Cell(nil, line)
	}

	pdf.SetFont("roboto", "", 18)
	pdf.SetX(pageWidth - marginLeft - 100)
	pdf.SetY(footerY + 15)
	pdf.Cell(nil, data.Total)

	// Footer note
	pdf.SetFont("roboto", "", 9)
	pdf.SetX(pageWidth - marginLeft - 180)
	pdf.SetY(pageWidth - 30)
	pdf.Cell(nil, "Send money abroad with Wise.")
}
