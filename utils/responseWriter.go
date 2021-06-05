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

func (writer *ResponseWriter) Write(data interface{}) error {
	switch writer.ctx.Get("accept", fiber.MIMEApplicationJSON) {
	case fiber.MIMEApplicationXML:
		result, err := xml.Marshal(data)
		if err != nil {
			return writer.ctx.Status(500).Send(writer.responseFail)
		}
		writer.ctx.Context().Response.Header.Add("Content-Type", fiber.MIMEApplicationXML)
		return writer.ctx.Send(result)
	default:
		return writer.ctx.JSON(data)
	}
}
