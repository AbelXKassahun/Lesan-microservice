package app

import (
	"context"
	"errors"
	"log"

	"lesson-service/internal/domain"
	"lesson-service/internal/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// type UserProgressService interface {
// 	GetUserProgress(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
// }

type UserProgressService struct {
	userProgressRepo ports.UserProgressRepository
	unitRepo         ports.UnitRepository
	sectionRepo      ports.SectionRepository
	lessonRepo       ports.LessonRepository
}

func NewUserProgressService(
	upr ports.UserProgressRepository,
	ur ports.UnitRepository,
	sr ports.SectionRepository,
	lr ports.LessonRepository,
) *UserProgressService {
	return &UserProgressService{
		userProgressRepo: upr,
		unitRepo:         ur,
		sectionRepo:      sr,
		lessonRepo:       lr,
	}
}

func (s *UserProgressService) GetUserProgress(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	// Fetch progress row
	progress, err := s.userProgressRepo.GetByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If no progress found, create a new one
			log.Println("No progress found so creating a new one")
			newProgress, createErr := s.CreateNewUserProgress(ctx, userID)
			if createErr != nil {
				return nil, createErr
			}
			if newProgress == nil {
				return nil, errors.New("failed to create new user progress: initial unit, section, or lesson not found")
			}
			progress = newProgress
		} else {
			return nil, err
		}
	}

	// Fetch units
	units, err := s.unitRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	var trail []map[string]interface{}

	for _, unit := range *units {
		sections, err := s.sectionRepo.GetByUnit(ctx, unit.ID)
		if err != nil {
			return nil, err
		}

		sectionList := []map[string]interface{}{}
		unitCompleted := true

		for _, section := range *sections {
			lessons, err := s.lessonRepo.GetBySection(ctx, section.ID)
			if err != nil {
				return nil, err
			}

			lessonList := []map[string]interface{}{}
			sectionCompleted := true

			for _, lesson := range *lessons {
				completed := lesson.ID != progress.CurrentLesson
				if !completed {
					sectionCompleted = false
					unitCompleted = false
				}
				lessonList = append(lessonList, map[string]interface{}{
					"lesson-id": lesson.ID,
					"completed": completed,
				})
			}

			sectionList = append(sectionList, map[string]interface{}{
				"id":           section.ID,
				"unit_id":      section.UnitID,
				"completed":    sectionCompleted,
				"lesson_count": section.LessonCount,
				"order_index":  section.OrderIndex,
				"lessons-id":   lessonList,
				"type":         section.Type,
			})
		}

		trail = append(trail, map[string]interface{}{
			"unit":          unit.OrderIndex,
			"unit-name":     unit.Title,
			"description":   unit.Description,
			"completed":     unitCompleted,
			"section_count": unit.SectionCount,
			"sections":      sectionList,
		})
	}

	result := map[string]interface{}{
		"user-progress": map[string]interface{}{
			"current-unit":       progress.CurrentUnit,
			"current_section_id": progress.CurrentSection,
			"next-lesson-id":     progress.CurrentLesson,
			"trail":              trail,
		},
	}

	return result, nil
}

func (s *UserProgressService) UpdateUserProgress(ctx context.Context, userID uuid.UUID) error {
	progress, err := s.userProgressRepo.GetByUser(ctx, userID)
	if err != nil {
		return err
	}

	lessons, err := s.lessonRepo.GetBySection(ctx, progress.CurrentSection)
	if err != nil {
		return err
	}

	currentLesson, err := s.lessonRepo.GetByID(ctx, progress.CurrentLesson)
	if err != nil {
		return err
	}
	// check if there are other lesson in the section
	nextLessonIndex := currentLesson.OrderIndex + 1
	for _, lesson := range *lessons {
		if nextLessonIndex == lesson.OrderIndex {
			progress.CurrentLesson = lesson.ID
			return s.userProgressRepo.Update(ctx, progress)
		}
	}

	//
	sections, err := s.sectionRepo.GetByUnit(ctx, progress.CurrentUnit)
	if err != nil {
		return err
	}

	currentSection, err := s.sectionRepo.GetByID(ctx, progress.CurrentSection)
	if err != nil {
		return err
	}

	nextSectionIndex := currentSection.OrderIndex + 1
	for _, section := range *sections {
		log.Println("section: ", section, " - order index: ", section.OrderIndex)
		if nextSectionIndex == section.OrderIndex {
			log.Println("matched section: ", section, " - matched order index: ", section.OrderIndex)
			progress.CurrentSection = section.ID
			lessons, err := s.lessonRepo.GetBySection(ctx, section.ID)
			if err != nil {
				return err
			}
			for _, lesson := range *lessons {
				if lesson.OrderIndex == 1 {
					progress.CurrentLesson = lesson.ID
					return s.userProgressRepo.Update(ctx, progress)
				}
			}
		}
	}

	currentUnit, err := s.unitRepo.GetByID(ctx, progress.CurrentUnit)
	if err != nil {
		return err
	}

	units, err := s.unitRepo.ListAll(ctx)
	if err != nil {
		return err
	}

	nextUnitIndex := currentUnit.OrderIndex + 1
	for _, unit := range *units {
		if nextUnitIndex == unit.OrderIndex {
			progress.CurrentUnit = unit.ID
			sections, err := s.sectionRepo.GetByUnit(ctx, unit.ID)
			if err != nil {
				return err
			}
			for _, section := range *sections {
				if section.OrderIndex == 1 {
					progress.CurrentSection = section.ID
					lessons, err := s.lessonRepo.GetBySection(ctx, section.ID)
					if err != nil {
						return err
					}
					for _, lesson := range *lessons {
						if lesson.OrderIndex == 1 {
							progress.CurrentLesson = lesson.ID
							return s.userProgressRepo.Update(ctx, progress)
						}
					}
				}
			}
		}
	}

	return nil
}

func (s *UserProgressService) CreateNewUserProgress(ctx context.Context, userID uuid.UUID) (*domain.UserProgress, error) {
	progress := &domain.UserProgress{
		UserID: userID,
	}

	units, err := s.unitRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, unit := range *units {
		if unit.OrderIndex == 1 {
			progress.CurrentUnit = unit.ID
			sections, err := s.sectionRepo.GetByUnit(ctx, unit.ID)
			if err != nil {
				return nil, err
			}
			for _, section := range *sections {
				if section.OrderIndex == 1 {
					progress.CurrentSection = section.ID
					lessons, err := s.lessonRepo.GetBySection(ctx, section.ID)
					if err != nil {
						return nil, err
					}
					for _, lesson := range *lessons {
						if lesson.OrderIndex == 1 {
							progress.CurrentLesson = lesson.ID
							return progress, s.userProgressRepo.Create(ctx, progress)
						}
					}
				}
			}

		}
	}
	return progress, nil
}
