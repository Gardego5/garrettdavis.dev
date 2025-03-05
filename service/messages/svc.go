package messages

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gardego5/garrettdavis.dev/model"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	db              *sqlx.DB
	listMessageAsc  *sqlx.NamedStmt
	listMessageDesc *sqlx.NamedStmt
	deleteMessage   *sqlx.NamedStmt
	createMessage   *sqlx.NamedStmt
}

func New(db *sqlx.DB) (*Service, error) {
	svc, err := Service{db: db}, error(nil)

	const LIST_TEMPLATE = `
SELECT ROW_NUMBER() OVER (ORDER BY created_at %s) as id,
       *
  FROM contact_messages
 LIMIT :limit
OFFSET :offset
`

	if svc.listMessageAsc, err = svc.db.PrepareNamed(
		fmt.Sprintf(LIST_TEMPLATE, "ASC")); err != nil {
		return nil, err
	}

	if svc.listMessageDesc, err = svc.db.PrepareNamed(
		fmt.Sprintf(LIST_TEMPLATE, "DESC")); err != nil {
		return nil, err
	}

	if svc.deleteMessage, err = svc.db.PrepareNamed(`
DELETE FROM contact_messages
WHERE id = :id
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

type (
	ListMessageInputSort string
	ListMessageInput     struct {
		Limit  int `db:"limit"`
		Offset int `db:"offset"`
		Sort   ListMessageInputSort
	}
)

const ListMessageInputSortASC ListMessageInputSort = "ASC"
const ListMessageInputSortDESC ListMessageInputSort = "DESC"

func (svc *Service) ListMessages(ctx context.Context, input *ListMessageInput) (out []model.ContactMessage, err error) {
	switch input.Sort {
	case ListMessageInputSortASC:
		err = svc.listMessageAsc.SelectContext(ctx, &out, input)
	case ListMessageInputSortDESC:
		err = svc.listMessageDesc.SelectContext(ctx, &out, input)
	}
	return
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

func (svc *Service) CountMessages(ctx context.Context) (count int, err error) {
	err = svc.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM contact_messages")
	return
}

func (svc *Service) Close() error {
	return errors.Join(
		svc.listMessageAsc.Close(),
		svc.listMessageDesc.Close(),
		svc.deleteMessage.Close(),
		svc.createMessage.Close(),
	)
}
