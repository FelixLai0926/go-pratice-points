package infrastructure

import (
	"fmt"
	"points/internal/domain/port"
	"points/internal/domain/valueobject"

	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
)

var _ port.Copier = (*CopierImpl)(nil)

type CopierImpl struct {
	option copier.Option
}

func NewCopierImpl() port.Copier {
	var opt copier.Option
	opt.Converters = append(opt.Converters, copier.TypeConverter{
		SrcType: decimal.Decimal{},
		DstType: valueobject.Money{},
		Fn: func(src interface{}) (interface{}, error) {
			if d, ok := src.(decimal.Decimal); ok {
				return valueobject.NewMoneyFromDecimal(d), nil
			}
			return nil, fmt.Errorf("cannot convert %T to Money", src)
		},
	})

	opt.Converters = append(opt.Converters, copier.TypeConverter{
		SrcType: valueobject.Money{},
		DstType: decimal.Decimal{},
		Fn: func(src interface{}) (interface{}, error) {
			if m, ok := src.(valueobject.Money); ok {
				return m.Value(), nil
			}
			return nil, fmt.Errorf("cannot convert %T to decimal.Decimal", src)
		},
	})
	return &CopierImpl{option: opt}
}

func (c *CopierImpl) Copy(dst, src interface{}) error {
	return copier.CopyWithOption(dst, src, c.option)
}
