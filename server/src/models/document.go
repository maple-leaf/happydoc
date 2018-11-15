package models

import (
	"errors"

	"github.com/thoas/go-funk"

	"github.com/lib/pq"

	"github.com/jinzhu/gorm"
)

type Document struct {
	gorm.Model
	Project  string `gorm:"unique"`
	Types    []DocumentType
	CVersion string `gorm:"-"`
	CType    string `gorm:"-"`
}

func (doc *Document) CreateOrUpdate(db *gorm.DB) error {
	if doc.CType == "" || doc.CVersion == "" {
		return errors.New("CType or CVersion is empty")
	}

	db.Where(doc).FirstOrCreate(&doc)

	docTypes := []DocumentType{}
	db.Model(&doc).Related(&docTypes)
	var docType DocumentType
	dt := funk.Find(docTypes, func(d DocumentType) bool {
		return d.Type == doc.CType
	})
	if dt == nil {
		docType = DocumentType{
			Type:       doc.CType,
			DocumentID: doc.ID,
		}
	} else {
		docType = dt.(DocumentType)
	}
	err := docType.CreateOrUpdate(db, doc.CVersion)

	return err
}

type DocumentType struct {
	gorm.Model
	Type       string         `gorm:"type:varchar(320)"`
	Versions   pq.StringArray `gorm:"type:varchar(16)[]"`
	DocumentID uint
}

func (docType *DocumentType) CreateOrUpdate(db *gorm.DB, version string) error {
	if version == "" {
		return errors.New("version is empty")
	}
	db.Where(docType).FirstOrCreate(&docType)
	modelForUpdate := DocumentType{}
	if !funk.ContainsString(docType.Versions, version) {
		modelForUpdate.Versions = append(docType.Versions, version)
	}

	db.Model(&docType).Update(modelForUpdate)

	return nil
}
