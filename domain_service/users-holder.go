package domain_service

import (
	"github.com/mstepan/user-service-golang/api"
	"github.com/mstepan/user-service-golang/domain"
	"math/rand"
	"sync"
)

type UserHolder struct {
	mutex sync.Mutex
	users map[string]*domain.UserProfile
}

func NewUserHolder() *UserHolder {
	return &UserHolder{users: make(map[string]*domain.UserProfile)}
}

func (ptr *UserHolder) AddUser(userProfile *api.CreateUserRequest) bool {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	if _, exists := ptr.users[userProfile.Username]; exists {
		return false
	}

	ptr.users[userProfile.Username] = &domain.UserProfile{Id: rand.Intn(10000), Username: userProfile.Username}

	return true
}

func (ptr *UserHolder) GetUserByUsername(username string) *domain.UserProfile {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	return ptr.users[username]
}
