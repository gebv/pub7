package admin

import (
	tarantool "crabler1/go-tarantool"
	"net/http"
	"strconv"

	"github.com/gebv/ff_tgbot/chats"
	"github.com/labstack/echo"
)

func DashboardHandler(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "dashboard", nil)
}

func ListChatsHandler(ctx echo.Context) error {
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	offset, _ := strconv.Atoi(ctx.QueryParam("offset"))
	list := listChats(offset, limit)
	return ctx.JSON(http.StatusOK, list)
}

func listChats(
	offset int,
	limit int,
) []chats.Chat {
	if limit <= 0 {
		limit = 100
	}
	limit = 1000
	var tuples []chats.Chat
	err := TarantoolConnection.SelectTyped(
		"ff_tgbot_statechats",
		"primary",
		uint32(offset),
		uint32(limit),
		tarantool.IterReq,
		[]interface{}{},
		&tuples,
	)
	if err != nil {
		log.Errorw(
			"get list states of chat",
			"err", err,
		)
		return []chats.Chat{}
	}
	if len(tuples) == 0 {
		return []chats.Chat{}
	}
	return tuples
}
