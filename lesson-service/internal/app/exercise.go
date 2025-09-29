package app

import(
	"context"

	"github.com/google/uuid"

	"lesson-service/internal/domain"
	"lesson-service/internal/ports"
)

type ExerciseService struct {
	Repo ports.ExerciseRepository
}

func NewExerciseService (repo ports.ExerciseRepository) *ExerciseService {	
	return &ExerciseService{Repo: repo}
}

func (s *ExerciseService) CreateExercise(ctx context.Context, exercise *domain.Exercise) (*domain.Exercise, error) {
	if exercise.ID == uuid.Nil {
		exercise.ID = uuid.New()
	}
	err := s.Repo.Create(ctx, exercise)
	return exercise, err
}

func (s *ExerciseService) GetExercise(ctx context.Context, id uuid.UUID) (*domain.Exercise, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *ExerciseService) GetExercisesByLesson(ctx context.Context, lessonID uuid.UUID) (*[]domain.Exercise, error) {
	return s.Repo.GetByLesson(ctx, lessonID)
}

func (s *ExerciseService) DeleteExercise(ctx context.Context, id uuid.UUID) error {
	return s.Repo.Delete(ctx, id)
}