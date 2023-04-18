package data

type UserPermissions []string

type Permissions struct {
	ID   uint64 `json:"id" gorm:"primaryKey"`
	Code string `json:"code" gorm:"unique;not null"`
}

func (p UserPermissions) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}
	return false
}
