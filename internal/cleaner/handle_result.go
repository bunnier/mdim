package cleaner

import "fmt"

type HandleResult struct {
	ImagePath string
	Err       error
}

func (result HandleResult) String() string {
	switch {
	case result.Err != nil:
		return fmt.Sprintf("[image handle]:Find a no reference image, but fail to delete.\n----> %s\n----> %s", result.ImagePath, result.Err.Error())
	default:
		return fmt.Sprintf("[image handle]:Delete a no reference image successfully.\n----> %s", result.ImagePath)
	}
}
