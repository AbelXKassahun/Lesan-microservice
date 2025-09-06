package app

import (
	"context"

	"github.com/google/uuid"
	"lesson-service/internal/ports"
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
		return nil, err
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
			"current-unit":      progress.CurrentUnit,
			"current_section_id": progress.CurrentSection,
			"next-lesson-id":    progress.CurrentLesson,
			"trail":             trail,
		},
	}

	return result, nil
}
