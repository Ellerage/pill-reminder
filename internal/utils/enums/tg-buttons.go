package enums

type Actions string

const (
	ActionTake   Actions = "Take"
	ActionEdit   Actions = "Edit"
	ActionCreate Actions = "/start"
)

type SendMessageButtons struct {
	Take bool
	Edit bool
}
