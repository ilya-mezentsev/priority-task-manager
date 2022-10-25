package account

import (
	log "github.com/sirupsen/logrus"
	"priority-task-manager/shared/pkg/repositories"
	"priority-task-manager/shared/pkg/types"
)

type Service struct {
	repository repositories.Reader[types.Account, string]
}

func MakeService(repository repositories.Reader[types.Account, string]) Service {
	return Service{
		repository: repository,
	}
}

func (s Service) GetAccount(hash string) (types.Account, error) {
	account, err := s.repository.Get(hash)
	if err != nil {
		log.WithFields(log.Fields{
			"hash": hash,
		}).Errorf("Unable to fetch account: %v", err)

		return types.Account{}, err
	}

	return account, nil
}
