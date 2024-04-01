package tm

import "fmt"

type ErrDuplicteMessage string

func (e ErrDuplicteMessage) Error() string {
	return fmt.Sprintf("duplicate message id encountered: %s", string(e))
}
