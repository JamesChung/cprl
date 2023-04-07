package client_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/smithy-go"

	"github.com/JamesChung/cprl/pkg/client"
)

func Test(t *testing.T) {
	c, err := client.NewSTSClient("default")
	if err != nil {
		fmt.Println(err)
	}
	out, err := c.GetCallerIdentity()
	if err != nil {
		var gen *smithy.GenericAPIError
		if errors.As(err, &gen) {
			fmt.Println(gen.Code)
			fmt.Println(gen.Message)
		}
		fmt.Println(err)
	}
	fmt.Println(out.Arn)
}
