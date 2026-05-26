package web_fs_repository

import (
	"fmt"
	"os"

	core_errors "github.com/swoggoi/todo_list/internal/core/errors"
)

func (r *WebRepository) GetFile(filePath string) ([]byte, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file : %s : %w", file, core_errors.ErrNotFound)
		}
		return nil, fmt.Errorf("get file : %s : %w", file, err)
	}

	return file, nil
}
