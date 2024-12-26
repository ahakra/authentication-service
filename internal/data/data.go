package data

type UserRepositoryInterface interface {
	Insert(user *UserModel) (*UserModel, error)
	GetByEmail(email string) (*UserModel, error)
	Update(user *UserModel) error
	GetById(userID int64) (*UserModel, error)
	UpdateUserActivationStatus(userID int64, status bool) error
}

type TokenRepositoryInterface interface {
	Insert(token *Token) (*Token, error)
	Delete(hash []byte) error
	GetByUserID(userID int64) ([]Token, error)
	GetByUserIDAndScope(userID int64, scope TokenScope) ([]Token, error)
	DeleteTokensForUser(userID int64, scope TokenScope) error
}

type PermissionsRepositoryInterface interface {
	InsertPermissions(permission string) error
	InsertUserPermissions(userID, permissionID int64) error
	DeleteUserPermissions(userID, permissionID int64) error
	GetPermissionIDByName(permission string) (int64, error)
	GetAllForUser(userID int64) (Permissions, error)
}
type RepoManager struct {
	UserRepo        UserRepositoryInterface
	TokenRepo       TokenRepositoryInterface
	PermissionsRepo PermissionsRepositoryInterface
}

// NewRepoManager creates a new instance of RepoManager with the given UserRepository
func NewRepoManager(userRepo UserRepositoryInterface, tokenRepo TokenRepositoryInterface, permissionRepo PermissionsRepositoryInterface) *RepoManager {
	return &RepoManager{
		UserRepo:        userRepo,
		TokenRepo:       tokenRepo,
		PermissionsRepo: permissionRepo,
	}
}
