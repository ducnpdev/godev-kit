package models

type TranslationModel struct {
	Source      string `json:"source" gorm:"column:source"`
	Destination string `json:"destination" gorm:"column:destination"`
	Original    string `json:"original" gorm:"column:original"`
	Translation string `json:"translation" gorm:"column:translation"`
}
