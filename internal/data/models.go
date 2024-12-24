package data

type UserRepositoryInterface interface {
	Insert(user *UserModel) (*UserModel, error)
	GetByEmail(email string) (*UserModel, error)
	Update(user *UserModel) error
}

type TokenRepositoryInterface interface {
	Insert(token *Token) (*Token, error)
	Delete(hash []byte) error
	GetByUserID(userID int64) ([]Token, error)
	GetByUserIDAndScope(userID int64, scope string) ([]Token, error)
}

type RepoManager struct {
	UserRepo  UserRepositoryInterface
	TokenRepo TokenRepositoryInterface
}

// NewRepoManager creates a new instance of RepoManager with the given UserRepository
func NewRepoManager(userRepo UserRepositoryInterface, tokenRepo TokenRepositoryInterface) *RepoManager {
	return &RepoManager{
		UserRepo:  userRepo,
		TokenRepo: tokenRepo,
	}
}
