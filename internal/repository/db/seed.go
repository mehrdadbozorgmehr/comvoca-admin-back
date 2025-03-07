package db

import (
	"fmt"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/entity"
	"gorm.io/gorm"
)

func SeedSpecialities(db *gorm.DB) error {
	specialities := map[string][]string{
		"General Dentistry": {
			"Dental Cleaning", "Dental Examination", "Dentail X-rays", "Deep Cleaning",
			"Dental Filings", "Root scaling and planing", "Fluoride Treatments", "Dental Sealants",
			"Cosmetic Dentistry", "Oral Surgery",
		},
		"Orthodontics Office": {
			"Diagnosis and Assessment", "Oral Examination", "Dental X-rays", "Diagnostic Models",
		},
		"Paediatrics Office": {
			"Gum Disease Treatment", "Scaling and Root Planing", "Soft Tissue Surgery",
			"Laser Therapy", "Bone Grafting", "Dental Implant Placement", "Implant Surgery",
			"Bone Grafting (for implant placement)", "Cosmetic Periodental Procedures",
			"Gum Contouring", "Gummy Smile Correction",
		},
		"Endodontics Office": {
			"Root Canal Therapy", "Root Canal Treatment", "Retreatment", "Apicoectomy",
			"Trauma Treatment", "Treatment of Dental Injuries",
		},
		"Prosthodontics Offices": {
			"Dental Implants", "Implant placement surgery", "Bone grating to prepare for implant placement",
			"Abutment placement", "Crown placement", "Dentures", "Crown and Bridges",
		},
		"Oral and Maxillofacial Surgery": {
			"Dental Procedures", "Wisdom Tooth Extraction", "Dental Implant Surgery", "Bone Grafting",
			"Soft Tissue Grafting", "Tooth Extractions", "Facial Trauma Surgery", "Facial Fracture Repair",
			"Soft Tissue Repair", "Orthognathic Surgery", "Jaw Surgery", "Cleft Lip and Palate Repair",
			"Cosmetic Procedures", "Facial Cosmetic Surgery", "Dental Implants for Cosmetic Purposes",
		},
	}

	for parentName, childNames := range specialities {
		parent := entity.Speciality{Name: parentName}
		if err := db.Create(&parent).Error; err != nil {
			return fmt.Errorf("failed to insert parent %s: %w", parentName, err)
		}

		for _, childName := range childNames {
			child := entity.Speciality{Name: childName, ParentID: &parent.ID}
			if err := db.Create(&child).Error; err != nil {
				return fmt.Errorf("failed to insert child %s: %w", childName, err)
			}
		}
	}

	return nil
}
