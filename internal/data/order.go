package data

type Order struct {
	ID              uint64       `json:"id" gorm:"primaryKey"`
	SenderID        uint64       `json:"sender_id"`
	SenderAddress   string       `json:"sender_address"`
	ReceiverID      uint64       `json:"receiver_id"`
	ReceiverAddress string       `json:"receiver_address"`
	Mass            float32      `json:"mass" sql:"type:decimal(10,2);"`
	InsuranceID     uint64       `json:"insurance_id"`
	Value           float32      `json:"value" sql:"type:decimal(10,2);"`
	Coverage        float32      `json:"coverage" sql:"type:decimal(3,2);"`
	Properties      []Properties `json:"-" gorm:"many2many:order_properties;"`
}

type Properties struct {
	ID   uint16 `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}
