package valueobject

import (
	"points/internal/shared/apperror"
	"points/internal/shared/errcode"

	"github.com/shopspring/decimal"
)

type Money struct {
	value decimal.Decimal
}

var Zero Money = Money{value: decimal.NewFromInt(0)}

func NewMoneyFromString(s string) (Money, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return Money{}, apperror.Wrap(errcode.ErrInvalidRequest, "invalid money value", err)
	}
	return Money{value: d}, nil
}

func NewMoneyFromDecimal(d decimal.Decimal) Money {
	return Money{value: d}
}

func (m Money) Add(other Money) Money {
	return Money{value: m.value.Add(other.value)}
}
func (m Money) Sub(other Money) Money {
	return Money{value: m.value.Sub(other.value)}
}

func (m Money) Subtract(other Money) Money {
	return Money{value: m.value.Sub(other.value)}
}

func (m Money) Multiply(factor decimal.Decimal) Money {
	return Money{value: m.value.Mul(factor)}
}

func (m Money) Equals(other Money) bool {
	return m.value.Equal(other.value)
}

func (m Money) Value() decimal.Decimal {
	return m.value
}

func (m Money) String() string {
	return m.value.String()
}

func (m Money) Cmp(other Money) int {
	return m.value.Cmp(other.value)
}

func (m Money) LessThan(other Money) bool {
	return m.value.LessThan(other.value)
}

func (m Money) GreaterThan(other Money) bool {
	return m.value.GreaterThan(other.value)
}

func (m Money) MarshalJSON() ([]byte, error) {
	return m.value.MarshalJSON()
}
