package data

type Insurance struct {
	ID          uint64  `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name" gorm:"unique;not null;size:255"`
	Description string  `json:"description" gorm:"not null"`
	Region      string  `json:"region" gorm:"size:255"`
	Orders      []Order `json:"-" gorm:"foreignKey:InsuranceID;references:ID"`
}
