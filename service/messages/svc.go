package messages

import (
	"context"
	"errors"

	"github.com/Gardego5/garrettdavis.dev/model"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	db            *sqlx.DB
	listMessage   *sqlx.NamedStmt
	deleteMessage *sqlx.NamedStmt
	createMessage *sqlx.NamedStmt
}

func New(db *sqlx.DB) (*Service, error) {
	svc, err := Service{db: db}, error(nil)

	if svc.listMessage, err = svc.db.PrepareNamed(`
SELECT ROW_NUMBER() OVER (ORDER BY created_at) as id,
       *
  FROM contact_messages
 LIMIT :limit
OFFSET :offset
`); err != nil {
		return nil, err
	}

	if svc.deleteMessage, err = svc.db.PrepareNamed(`
DELETE FROM contact_messages
 WHERE id = ?
`); err != nil {
		return nil, err
	}

	if svc.createMessage, err = svc.db.PrepareNamed(`
INSERT INTO contact_messages
     ( name
     , email
     , message
     , created_at)
VALUES (:name, :email, :message, :created_at)`); err != nil {
		return nil, err
	}

	return &svc, nil
}

type ListMessageInput struct {
	Limit  int `db:"limit"`
	Offset int `db:"offset"`
}

func (svc *Service) ListMessages(ctx context.Context, input *ListMessageInput) ([]model.ContactMessage, error) {
	output := []model.ContactMessage{}
	err := svc.listMessage.SelectContext(ctx, output, input)
	return output, err
}

type DeleteMessageInput struct {
	ID int `db:"id"`
}

func (svc *Service) DeleteMessage(ctx context.Context, input *DeleteMessageInput) error {
	_, err := svc.deleteMessage.ExecContext(ctx, input)
	return err
}

func (svc *Service) CreateMessage(ctx context.Context, input *model.ContactMessage) error {
	_, err := svc.createMessage.ExecContext(ctx, input)
	return err
}

func (svc *Service) Close() error {
	return errors.Join(
		svc.listMessage.Close(),
		svc.deleteMessage.Close(),
		svc.createMessage.Close(),
	)
}
