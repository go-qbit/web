package form

import "context"

type IAntiCSRFConfig interface {
	GetTokenSalt() string
	GetTokenLength() int
	GetErrorText() string
	GetCtxParser() func(ctx context.Context) (userID uint32, formPath string)
}

type AntiCSRFConfig struct {
	Salt      string
	Length    int
	ErrorText string
	CtxParser func(ctx context.Context) (userID uint32, formPath string)
}

func (c *AntiCSRFConfig) GetTokenSalt() string {
	return c.Salt
}

func (c *AntiCSRFConfig) GetTokenLength() int {
	return c.Length
}

func (c *AntiCSRFConfig) GetErrorText() string {
	return c.ErrorText
}

func (c *AntiCSRFConfig) GetCtxParser() func(ctx context.Context) (userID uint32, formPath string) {
	return c.CtxParser
}
