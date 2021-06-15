package validators

import (
	"context"
	"fmt"

	"github.com/smancke/mailck"
)

const fromEmail = "moscow.beverage@gmail.com"

var ErrLongEmailLen = fmt.Errorf("length of email is too long")
var ErrInvalidEmailFormat = fmt.Errorf("invalid email format")

func ValidateEmail(ctx context.Context, email string) error {
	result, err := mailck.CheckWithContext(ctx, fromEmail, email)
	if err != nil {
		return err
	}
	if !result.IsValid() {
		return fmt.Errorf(result.Message)
	}
	return nil
}
