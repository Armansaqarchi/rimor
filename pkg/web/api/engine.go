package api

import (
	"rimor/pkg/engine"

	"github.com/gofiber/fiber/v2"
)




type EngineHandler struct {
	Engine engine.Engine
}

var e *EngineHandler

func GetEngineHandler() *EngineHandler{
	if e != nil {
		return e
	}

	return newEngineHandler()


}


func newEngineHandler() *EngineHandler{
	return &EngineHandler{
		// TODO
	}
}


func (e EngineHandler) Query(ctx *fiber.Ctx) error{
	if text := ctx.Query("text", ""); text == "" {
		return ctx.Status(400).JSON(
			map[string]any{
				"error": "text is not present",
			},
		)
	} else {
		res , err := e.Engine.Query(text)
		if err != nil {
			return err
		}

		return ctx.Status(200).JSON(
			map[string]any{
				"body": res.DocList,
			},
		)
	}

	
}