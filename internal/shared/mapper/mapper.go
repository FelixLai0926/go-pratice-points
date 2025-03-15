package mapper

import (
	"points/internal/domain/port"
	"points/internal/shared/apperror"
	"points/internal/shared/errcode"
)

func MapStruct[TDestination any, TSource any](config port.Config, ormAccount *TSource) (*TDestination, error) {
	entityDestination := new(TDestination)
	if err := config.Copy(entityDestination, ormAccount); err != nil {
		return nil, apperror.Wrap(errcode.ErrInternal, "failed to copy struct", err)
	}

	return entityDestination, nil
}
