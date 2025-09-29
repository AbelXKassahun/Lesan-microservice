package ports

import (
	"context"
	"lesson-service/internal/domain"

	"github.com/google/uuid"
)

type UserProgressRepository interface {
	Create(ctx context.Context, userProgress *domain.UserProgress) error
	GetByUser(ctx context.Context, userID uuid.UUID) (*domain.UserProgress, error)
	Update(ctx context.Context, userProgress *domain.UserProgress) error
}

type UnitRepository interface {
	Create(ctx context.Context, unit *domain.Unit) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Unit, error)
	ListAll(ctx context.Context) (*[]domain.Unit, error)
	Update(ctx context.Context, unit *domain.Unit) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type SectionRepository interface {
	Create(ctx context.Context, section *domain.Section) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Section, error)
	GetByUnit(ctx context.Context, unitID uuid.UUID) (*[]domain.Section, error)
	ListAll(ctx context.Context) (*[]domain.Section, error)
	Update(ctx context.Context, section *domain.Section) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type LessonRepository interface {
	Create(ctx context.Context, lesson *domain.Lesson) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Lesson, error)
	GetByUnit(ctx context.Context, unitID uuid.UUID) (*[]domain.Lesson, error)
	GetBySection(ctx context.Context, sectionID uuid.UUID) (*[]domain.Lesson, error)
	ListAll(ctx context.Context) (*[]domain.Lesson, error)
	Update(ctx context.Context, lesson *domain.Lesson) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ExerciseRepository interface {
	Create(ctx context.Context, e *domain.Exercise) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Exercise, error)
	GetByLesson(ctx context.Context, lessonID uuid.UUID) (*[]domain.Exercise, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
