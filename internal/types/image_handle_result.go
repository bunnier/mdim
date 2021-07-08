package types

import "fmt"

type ImageHandleResult struct {
	ImagePath string
	Deleted   bool
	Err       error
}

func (result ImageHandleResult) ToString() string {
	switch {
	case result.Err != nil:
		return fmt.Sprintf("[image handle]:Find a no reference image, but fail to delete.\n----> %s\n----> %s", result.ImagePath, result.Err.Error())
	case result.Deleted:
		return fmt.Sprintf("[image handle]:Delete a no reference image successfully.\n----> %s", result.ImagePath)
	case !result.Deleted:
		return fmt.Sprintf("[image handle]:Find a no reference image, do not delete this time.\n----> %s", result.ImagePath)
	default:
		return "Impossible error."
	}
}
