package enums

type Actions string

const (
	ActionTake   Actions = "Take"
	ActionEdit   Actions = "Edit"
	ActionDelay  Actions = "Delay"
	ActionCreate Actions = "/start"
)

type SendMessageButtons struct {
	Take  bool
	Edit  bool
	Delay bool
}
