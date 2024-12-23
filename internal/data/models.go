package data

type UserRepositoryInterface interface {
	Insert(user *UserModel) (*UserModel, error)
	GetByEmail(email string) (*UserModel, error)
	Update(user *UserModel) error
}

type RepoManager struct {
	UserRepo UserRepositoryInterface
}
type RepoManagerInterface interface {
	UserRepositoryInterface
}

// NewRepoManager creates a new instance of RepoManager with the given UserRepository
func NewRepoManager(userRepo UserRepositoryInterface) *RepoManager {
	return &RepoManager{UserRepo: userRepo}
}
