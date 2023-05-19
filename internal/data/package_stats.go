package data

type PackageStats struct {
	ID                uint64  `json:"id" gorm:"primaryKey"`
	Mean              float64 `json:"mean" gorm:"not null"`
	Variance          float64 `json:"variance" gorm:"not null"`
	StandardDeviation float64 `json:"standard_deviation" gorm:"not null"`
	Prediction        float64 `json:"prediction" gorm:"not null"`
}
