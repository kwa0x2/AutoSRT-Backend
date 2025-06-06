package usecase

import (
	"context"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"time"
)

type contactUseCase struct {
	contactBaseRepository domain.BaseRepository[*domain.Contact]
}

func NewContactUseCase(contactBaseRepository domain.BaseRepository[*domain.Contact]) domain.ContactUseCase {
	return &contactUseCase{contactBaseRepository: contactBaseRepository}
}

func (cr *contactUseCase) Create(contact *domain.Contact) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now().UTC()
	contact.CreatedAt = now
	contact.UpdatedAt = now

	return cr.contactBaseRepository.Create(ctx, contact)
}
