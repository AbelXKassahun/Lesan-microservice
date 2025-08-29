package domain

import (
	"time"
)

// UserBadge represents badges earned by a user
type UserBadge struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	UserID    string `gorm:"index;type:text"`
	BadgeName string
	AwardedAt time.Time
}

type BadgeNames int

const(
	FirstLesson BadgeNames = iota
    BackForMore
	FirstWeek
)

func (badge BadgeNames) String() string {
    switch badge {
    case FirstLesson:
        return "FirstLesson"
    case BackForMore:
        return "BackForMore"
    case FirstWeek:
        return "FirstWeek"
    default:
        return ""
    }
}
// func (badge BadgeNames) String() (string, error) {
//     switch badge {
//     case FirstLesson:
//         return "FirstLesson", nil
//     case BackForMore:
//         return "BackForMore", nil
//     case FirstWeek:
//         return "FirstWeek", nil
//     default:
//         return "", fmt.Errorf("invalid badge name")
//     }
// }

func IsValidStatus(badge BadgeNames) bool {
    switch badge {
    case FirstLesson, BackForMore, FirstWeek:
        return true
    default:
        return false
    }
}
