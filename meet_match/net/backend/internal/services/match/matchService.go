package match_service

import (
	"log"
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/services/cards/repository"
	match_repo "test_backend_frontend/internal/services/match/matchRepo"
	session "test_backend_frontend/internal/sessions"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type IMatchService interface {
	GetMatchedCardsBySession(sessionID uuid.UUID) ([]*models.Card, error)
}

func NewMatchService(matchRepoSrc match_repo.IMatchRepo, sessionManSc session.SessionManager, cardRepoSrc repository.CardRepository) IMatchService {
	return &MatchService{
		matchRepo:  matchRepoSrc,
		sessionMan: sessionManSc,
		cardRepo:   cardRepoSrc,
	}

}

type MatchService struct {
	matchRepo  match_repo.IMatchRepo
	sessionMan session.SessionManager
	cardRepo   repository.CardRepository
}

func (m *MatchService) GetMatchedCardsBySession(sessionID uuid.UUID) ([]*models.Card, error) {
	users, err := m.sessionMan.GetUsers(sessionID)
	if err != nil {
		return nil, errors.Wrap(err, "error in matchService retriving users from session")
	}
	matches, err := m.matchRepo.GetUserMatchesBySession(sessionID, users[0].ID) // note that the first user can be any user from session
	if err != nil {
		return nil, err
	}
	matchedCards := make([]*models.Card, len(matches))

	for i, match := range matches {
		matchedCards[i], err = m.cardRepo.GetCard(match.CardMatchedID)
		if err != nil {
			log.Printf("error extracting matchedCards by session: %w\n", err) //TODO:: add logging
		}

	}
	return matchedCards, nil
}
