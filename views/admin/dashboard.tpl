{{define "dashboard"}}
<!doctype html>
<meta charset="utf-8">
<title>Dashboard</title>
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/uikit/2.27.4/css/uikit.almost-flat.min.css"/>

<script src="https://cdnjs.cloudflare.com/ajax/libs/mithril/1.1.3/mithril.min.js"></script> 
<script src="https://cdnjs.cloudflare.com/ajax/libs/lodash.js/4.17.4/lodash.min.js"></script> 

<div id="root"></div>
<script>
var Dashboard = {
    limit: 50,
    offset: 0,
    list: [],
    oninit: function() {
        m.request({
            method: "GET",
            url: "/api/v1/chats/states?" +
                m.buildQueryString({offset: this.offset, limit: this.limit})
        }).then(function(res) {
            this.list = res;
        }.bind(this))
    },
    view: function(v) {
        var list = _.map(v.state.list, function(item) {
            return m("li", m(ChatInfo, {key: item.ChatID, data: item}))
        })
        return m(
            "div.uk-grid", 
            m(
                "div.uk-width-1-2 uk-container-center", 
                m("ul.uk-list uk-list-space", list),
            ),
        )
    },
};
var ChatInfo = {
    view: function(v) {
        var data = v.attrs.data;
        return m("div.uk-panel uk-panel-box uk-width-1-1", [
            m("div.uk-panel-badge", [
                m("span.uk-badge", data.ChatID), 
                ] ),
            m("div.uk-panel-title", data.Props["user_name"] || "not set username"),
            m("ul.uk-list", _.map(data.Props, function(item, key) {
                return m("li", [
                    m("span.uk-text-bold uk-margin-right", key),
                    m("span", item),
                ])
            })),
            m("hr"),
            m("ul.uk-text-small", [
                m("li", [
                    m("span.uk-text-bold uk-margin-right", "Последний # вопрос"),
                    m("span", data.LastQID),
                ]),
                m("li", [
                    m("span.uk-text-bold uk-margin-right", "Текущий # сценария"),
                    m("span", data.ScriptID),
                ]),
                m("li", [
                    m("span.uk-text-bold uk-margin-right", "Последняя дата обновления"),
                    m("span", data.UpdatedAt),
                ]),
            ])
        ])
    }
}
m.mount(document.getElementById("root"), Dashboard)
</script>

{{end}}