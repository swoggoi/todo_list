package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/swoggoi/todo_list/internal/core/errors"
)

type User struct {
	ID      int
	Version int

	FullName    string
	PhoneNumber *string // необязательно

}

func NewUser(
	id int,
	version int,
	fullName string,
	phoneNumber *string,
) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func NewUserUnitialized(
	fullName string,
	phoneNumber *string,
) User {
	return NewUser(
		UnitializedID,
		UnitializedVersion,
		fullName,
		phoneNumber,
	)
}

func NewUserPatch(
	fullName Nullable[string],
	phoneNumber Nullable[string],
) UserPatch {
	return UserPatch{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func (u *User) Validate() error {
	fullNameLenght := len([]rune(u.FullName))
	if fullNameLenght < 3 || fullNameLenght > 100 {
		return fmt.Errorf(
			"invalid `FullName` len: %d: %w",
			fullNameLenght,
			core_errors.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber != nil {
		phoneNumberLenght := len([]rune(*u.PhoneNumber))
		if phoneNumberLenght < 10 || phoneNumberLenght > 15 {
			return fmt.Errorf(
				"invalid PhoneNumber len: %d: %w",
				phoneNumberLenght,
				core_errors.ErrInvalidArgument,
			)
		}
		re := regexp.MustCompile(`^\+[0-9]+$`)
		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf(
				"invalid PhoneNumber len: %d: %w",
				phoneNumberLenght,
				core_errors.ErrInvalidArgument,
			)
		}

	}

	return nil
}

type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

func (p *UserPatch) Validate() error {
	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf(
			"'FullName' can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	return nil
}
func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u
	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}
	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp

	return nil

}
