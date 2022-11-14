package service

import (
	"github.com/google/uuid"
	"github.com/mstepan/user-service-golang/domain"
	"sync"
)

type UserHolder struct {
	mutex sync.Mutex
	users map[string]*domain.UserProfile
}

func NewUserHolder() *UserHolder {
	return &UserHolder{users: make(map[string]*domain.UserProfile)}
}

func (ptr *UserHolder) AddUser(username string) *domain.UserProfile {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	if _, exists := ptr.users[username]; exists {
		return nil
	}

	randomUuid := uuid.New()

	userProfile := &domain.UserProfile{Id: randomUuid.String(), Username: username}

	ptr.users[userProfile.Username] = userProfile

	return userProfile
}

func (ptr *UserHolder) GetUserByUsername(username string) *domain.UserProfile {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	return ptr.users[username]
}

func (ptr *UserHolder) GetAllUsers() []*domain.UserProfile {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	allUsers := make([]*domain.UserProfile, 0, len(ptr.users))

	for _, userProfile := range ptr.users {
		allUsers = append(allUsers, userProfile)
	}

	return allUsers
}

func (ptr *UserHolder) DeleteUserByUsername(username string) bool {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	if _, exists := ptr.users[username]; exists {
		delete(ptr.users, username)
		return true
	}

	return false
}

func (ptr *UserHolder) GetUsersCount() int {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	return len(ptr.users)
}
