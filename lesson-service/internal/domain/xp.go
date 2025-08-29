package domain

import (
	// "gorm.io/gorm"
)

// UserXP represents total XP a user has earned
type UserXP struct {
	UserID string `gorm:"primaryKey;type:text"`
	Total  int
	League League
}

type League string

const (
    Iron      League = "Iron"
    Bronze    League = "Bronze"
    Silver    League = "Silver"
    Gold      League = "Gold"
    Platinum  League = "Platinum"
    Diamond   League = "Diamond"
    Ascendant League = "Ascendant"
    Immortal  League = "Immortal"
    Radiant   League = "Radiant"
)

var leagueThresholds = []struct {
    League League
    MinXP  int
}{
    {Iron, 0},
    {Bronze, 1000},
    {Silver, 2500},
    {Gold, 5000},
    {Platinum, 8000},
    {Diamond, 12000},
    {Ascendant, 17000},
    {Immortal, 23000},
    {Radiant, 30000},
}

func DetermineLeague(totalXP int) League {
    var league League
    for _, lt := range leagueThresholds {
        if totalXP >= lt.MinXP {
            league = lt.League
        } else {
            break
        }
    }
    return league
}