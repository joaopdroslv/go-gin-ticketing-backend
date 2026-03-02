package accesscontrol

import (
	"context"
	models "go-gin-ticketing-backend/internal/access_control/models"
	repository "go-gin-ticketing-backend/internal/access_control/repository"
	schemas "go-gin-ticketing-backend/internal/access_control/schemas"
	"go-gin-ticketing-backend/internal/domain"
	"log"
)

type PermissionService struct {
	permissionRepository repository.PermissionRepository
}

func NewPermissionService(
	permissionRepository repository.PermissionRepository,
) *PermissionService {

	return &PermissionService{permissionRepository: permissionRepository}
}

// TODO: finish this method
func (s *PermissionService) GetAllPermissions(
	ctx context.Context,
	query *schemas.GetAllPermissionsQuery,
) {

	pagination := domain.NewPagination(query.Page, query.Limit)

	// NOTE: No need for a dto for now
	name := query.Name

	permissions, total, err := s.permissionRepository.GetAllPermissions(ctx, name, pagination)
	if err != nil {
		log.Println(err)
		// return nil, err
	}

	log.Println(permissions)
	log.Println(total)

	_ = permissions
	_ = total
}

func (s *PermissionService) GetPermissionsByRoleID(
	ctx context.Context,
	id int64,
) ([]models.Permission, error) {

	return s.permissionRepository.GetPermissionsByRoleID(ctx, id)
}

func (s *PermissionService) UserHasPermission(
	ctx context.Context,
	userID int64,
	requiredPermission string,
) (bool, error) {

	return s.permissionRepository.UserHasPermission(ctx, userID, requiredPermission)
}
