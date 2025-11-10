package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/srcndev/message-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	cleanup := func() {
		sqlDB.Close()
	}

	return db, mock, cleanup
}

func TestMessageRepository_Create_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewMessageRepository(db)

	message := &domain.Message{
		PhoneNumber: "+905551234567",
		Content:     "Test message",
		Status:      domain.StatusPending,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "messages"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), message)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_Create_Error(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewMessageRepository(db)

	message := &domain.Message{
		PhoneNumber: "+905551234567",
		Content:     "Test message",
		Status:      domain.StatusPending,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "messages"`)).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	err := repo.Create(context.Background(), message)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetByID_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewMessageRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"phone_number", "content", "status", "message_id", "sent_at",
	}).AddRow(
		1, now, now, nil,
		"+905551234567", "Test message", domain.StatusPending, nil, nil,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "messages" WHERE "messages"."id" = $1`)).
		WithArgs(1, 1).
		WillReturnRows(rows)

	message, err := repo.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, uint(1), message.ID)
	assert.Equal(t, "+905551234567", message.PhoneNumber)
	assert.Equal(t, "Test message", message.Content)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetByID_NotFound(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewMessageRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "messages" WHERE "messages"."id" = $1`)).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	message, err := repo.GetByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, message)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_List_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewMessageRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"phone_number", "content", "status", "message_id", "sent_at",
	}).
		AddRow(1, now, now, nil, "+905551111111", "Message 1", domain.StatusPending, nil, nil).
		AddRow(2, now, now, nil, "+905552222222", "Message 2", domain.StatusSent, "msg-id", &now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "messages"`)).
		WillReturnRows(rows)

	messages, err := repo.List(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.Len(t, messages, 2)
	if len(messages) >= 2 {
		assert.Equal(t, "+905551111111", messages[0].PhoneNumber)
		assert.Equal(t, "+905552222222", messages[1].PhoneNumber)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_List_Empty(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewMessageRepository(db)

	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"phone_number", "content", "status", "message_id", "sent_at",
	})

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "messages"`)).
		WillReturnRows(rows)

	messages, err := repo.List(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.Empty(t, messages)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetPendingMessages_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewMessageRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"phone_number", "content", "status", "message_id", "sent_at",
	}).
		AddRow(1, now, now, nil, "+905551111111", "Pending 1", domain.StatusPending, nil, nil).
		AddRow(2, now, now, nil, "+905552222222", "Pending 2", domain.StatusPending, nil, nil)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "messages" WHERE status = $1`)).
		WillReturnRows(rows)

	messages, err := repo.GetPendingMessages(context.Background(), 2)

	assert.NoError(t, err)
	assert.Len(t, messages, 2)
	if len(messages) >= 2 {
		assert.Equal(t, domain.StatusPending, messages[0].Status)
		assert.Equal(t, domain.StatusPending, messages[1].Status)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetPendingMessages_NoPending(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewMessageRepository(db)

	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"phone_number", "content", "status", "message_id", "sent_at",
	})

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "messages" WHERE status = $1`)).
		WillReturnRows(rows)

	messages, err := repo.GetPendingMessages(context.Background(), 2)

	assert.NoError(t, err)
	assert.Empty(t, messages)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetSentMessages_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewMessageRepository(db)

	sentAt := time.Now()
	msgID1 := "msg-123"
	msgID2 := "msg-456"

	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"phone_number", "content", "status", "message_id", "sent_at",
	}).
		AddRow(1, time.Now(), time.Now(), nil, "+905551234567", "Message 1", domain.StatusSent, msgID1, sentAt).
		AddRow(2, time.Now(), time.Now(), nil, "+905551234568", "Message 2", domain.StatusSent, msgID2, sentAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "messages" WHERE status = $1`)).
		WillReturnRows(rows)

	messages, err := repo.GetSentMessages(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.Len(t, messages, 2)
	if len(messages) == 2 {
		assert.Equal(t, domain.StatusSent, messages[0].Status)
		assert.Equal(t, domain.StatusSent, messages[1].Status)
		assert.NotNil(t, messages[0].MessageID)
		assert.NotNil(t, messages[1].MessageID)
		assert.Equal(t, msgID1, *messages[0].MessageID)
		assert.Equal(t, msgID2, *messages[1].MessageID)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetSentMessages_NoSent(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewMessageRepository(db)

	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"phone_number", "content", "status", "message_id", "sent_at",
	})

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "messages" WHERE status = $1`)).
		WillReturnRows(rows)

	messages, err := repo.GetSentMessages(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.Empty(t, messages)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_Update_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewMessageRepository(db)

	message := &domain.Message{
		ID:          1,
		PhoneNumber: "+905551234567",
		Content:     "Updated message",
		Status:      domain.StatusSent,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "messages"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), message)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_Update_Error(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewMessageRepository(db)

	message := &domain.Message{
		ID:          1,
		PhoneNumber: "+905551234567",
		Content:     "Updated message",
		Status:      domain.StatusSent,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "messages"`)).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	err := repo.Update(context.Background(), message)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_Delete_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewMessageRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "messages" SET "deleted_at"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_Delete_NotFound(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewMessageRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "messages" SET "deleted_at"`)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), 999)

	// GORM doesn't return error for soft delete even if not found
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_InterfaceCompliance(t *testing.T) {
	var _ MessageRepository = (*messageRepository)(nil)

	db, _, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewMessageRepository(db)
	assert.NotNil(t, repo)
}
