package repository

import (
	"github.com/google/uuid"
	"github.com/ngoctrng/bookz/internal/exchange"
	"gorm.io/gorm"
)

const tblProposals = "proposals"

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(p *exchange.Proposal) error {
	schema := DomainToProposalSchema(p)
	return r.db.Table(tblProposals).Create(schema).Error
}

func (r *Repository) GetByID(id int) (*exchange.Proposal, error) {
	var schema ProposalSchema
	err := r.db.Table(tblProposals).Where("id = ?", id).First(&schema).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return ProposalSchemaToDomain(&schema), nil
}

func (r *Repository) GetAll(uid uuid.UUID) ([]*exchange.Proposal, error) {
	var schemas []ProposalSchema
	err := r.db.Table(tblProposals).Where("request_by = ?", uid).Find(&schemas).Error
	if err != nil {
		return nil, err
	}

	proposals := make([]*exchange.Proposal, 0, len(schemas))
	for _, s := range schemas {
		proposals = append(proposals, ProposalSchemaToDomain(&s))
	}

	return proposals, nil
}
