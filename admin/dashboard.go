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
) []chats.State {
	if limit <= 0 {
		limit = 100
	}
	var tuples []chats.State
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
		return []chats.State{}
	}
	if len(tuples) == 0 {
		return []chats.State{}
	}
	return tuples
}
