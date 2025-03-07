package entity

type Speciality struct {
	ID              int          `json:"id"`
	Name            string       `json:"name"`
	ParentID        *int         `gorm:"index" json:"parent_id,omitempty"`
	Parent          *Speciality  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	SubSpecialities []Speciality `gorm:"foreignKey:ParentID" json:"subSpecialities,omitempty"`
}
