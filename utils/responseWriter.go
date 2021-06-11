package utils

import (
	"encoding/xml"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ResponseWriter struct {
	ctx          *fiber.Ctx
	logger       *zap.Logger
	responseFail []byte
}

func NewResponseWriter(ctx *fiber.Ctx, logger *zap.Logger, responseFail []byte) *ResponseWriter {
	if ctx == nil || logger == nil {
		panic("failed to create response writer, ctx and logger must not be nil")
	}
	return &ResponseWriter{
		ctx:          ctx,
		logger:       logger,
		responseFail: responseFail,
	}
}

func (writer *ResponseWriter) Status(status int) *ResponseWriter {
	writer.ctx.Status(status)
	return writer
}

func (writer *ResponseWriter) Write(data interface{}) error {
	switch writer.ctx.Get("accept", fiber.MIMEApplicationJSON) {
	case fiber.MIMEApplicationXML:
		result, err := xml.Marshal(data)
		if err != nil {
			writer.logger.Error("failed to create xml", zap.Error(err))
			return writer.ctx.Status(500).Send(writer.responseFail)
		}
		writer.ctx.Context().Response.Header.Add("Content-Type", fiber.MIMEApplicationXML)
		return writer.ctx.Send(result)
	default:
		return writer.ctx.JSON(data)
	}
}
