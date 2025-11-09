package seed

import (
	"context"
	"fmt"
	"time"

	"github.com/jaswdr/faker"
	"github.com/srcndev/message-service/internal/domain"
	"github.com/srcndev/message-service/pkg/logger"
	"gorm.io/gorm"
)

// Seeder handles database seeding
type Seeder struct {
	db    *gorm.DB
	faker faker.Faker
}

// NewSeeder creates a new seeder instance
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{
		db:    db,
		faker: faker.New(),
	}
}

// Run executes all seeders
func (s *Seeder) Run() error {
	ctx := context.Background()

	// Check if we already have data
	var count int64
	if err := s.db.WithContext(ctx).Model(&domain.Message{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count messages: %w", err)
	}

	// If data exists, skip seeding
	if count > 0 {
		logger.Info("Database already has %d messages, skipping seed", count)
		return nil
	}

	logger.Info("Seeding database with fake data...")

	// Seed pending messages
	if err := s.seedPendingMessages(ctx, 20); err != nil {
		return fmt.Errorf("failed to seed pending messages: %w", err)
	}

	// Seed sent messages (history)
	if err := s.seedSentMessages(ctx, 5); err != nil {
		return fmt.Errorf("failed to seed sent messages: %w", err)
	}

	logger.Info("Database seeded successfully: 20 pending, 5 sent")
	return nil
}

// seedPendingMessages creates fake pending messages
func (s *Seeder) seedPendingMessages(ctx context.Context, count int) error {
	messages := make([]domain.Message, count)

	messageTemplates := []string{
		"Hello %s! Your order #%s has been confirmed.",
		"Hi %s, your appointment is scheduled for %s.",
		"Dear %s, your verification code is: %s",
		"%s, your package has been shipped! Track: %s",
		"Welcome %s! Your account has been created successfully.",
		"Reminder: %s, your subscription expires on %s",
		"Hi %s! Special offer: %s",
		"%s, your payment of $%s has been received.",
		"Dear %s, your reservation #%s is confirmed.",
		"%s, your password reset code is: %s",
	}

	for i := 0; i < count; i++ {
		template := messageTemplates[i%len(messageTemplates)]
		name := s.faker.Person().FirstName()

		var content string
		switch i % len(messageTemplates) {
		case 0:
			content = fmt.Sprintf(template, name, s.faker.RandomStringWithLength(8))
		case 1:
			futureTime := time.Now().Add(48 * time.Hour)
			content = fmt.Sprintf(template, name, futureTime.Format("Jan 02, 15:04"))
		case 2:
			content = fmt.Sprintf(template, name, s.faker.RandomStringWithLength(6))
		case 3:
			content = fmt.Sprintf(template, name, s.faker.UUID().V4())
		case 4:
			content = fmt.Sprintf(template, name)
		case 5:
			futureDate := time.Now().Add(720 * time.Hour)
			content = fmt.Sprintf(template, name, futureDate.Format("2006-01-02"))
		case 6:
			content = fmt.Sprintf(template, name, fmt.Sprintf("%d%% discount on all products!", s.faker.IntBetween(10, 50)))
		case 7:
			content = fmt.Sprintf(template, name, fmt.Sprintf("%.2f", s.faker.Float64(2, 50, 500)))
		case 8:
			content = fmt.Sprintf(template, name, s.faker.RandomStringWithLength(10))
		case 9:
			content = fmt.Sprintf(template, name, s.faker.RandomStringWithLength(6))
		}

		messages[i] = domain.Message{
			PhoneNumber: s.generateTurkishPhone(),
			Content:     content,
			Status:      domain.StatusPending,
		}
	}

	return s.db.WithContext(ctx).Create(&messages).Error
}

// seedSentMessages creates fake sent messages (for history)
func (s *Seeder) seedSentMessages(ctx context.Context, count int) error {
	messages := make([]domain.Message, count)

	for i := 0; i < count; i++ {
		sentAt := time.Now().Add(-time.Duration(s.faker.IntBetween(1, 72)) * time.Hour)
		messageID := s.faker.UUID().V4()

		messages[i] = domain.Message{
			PhoneNumber: s.generateTurkishPhone(),
			Content:     fmt.Sprintf("This is a sent message: %s", s.faker.Lorem().Sentence(10)),
			Status:      domain.StatusSent,
			MessageID:   &messageID,
			SentAt:      &sentAt,
		}
	}

	return s.db.WithContext(ctx).Create(&messages).Error
}

// generateTurkishPhone generates a realistic Turkish phone number
func (s *Seeder) generateTurkishPhone() string {
	// Turkish mobile operators: 50x, 51x, 52x, 53x, 54x, 55x
	operators := []string{"505", "506", "507", "530", "531", "532", "533", "534", "535", "536", "537", "538", "539", "541", "542", "543", "544", "545", "546", "547", "548", "549", "551", "552", "553", "554", "555", "559"}
	operator := operators[s.faker.IntBetween(0, len(operators)-1)]

	// Generate 7 random digits
	number := s.faker.IntBetween(1000000, 9999999)

	return fmt.Sprintf("+90%s%d", operator, number)
}
