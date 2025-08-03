package enums

type Actions string
type ActionCommands string

const (
	ActionTake      Actions = "Take"
	ActionEdit      Actions = "Edit"
	ActionDelay     Actions = "Delay"
	ActionCreate    Actions = "/start"
	ActionMySetting Actions = "Settings"
	ActionUndo      Actions = "Undo"
)

const (
	ActionCommandTake     ActionCommands = "/take"
	ActionCommandEdit     ActionCommands = "/edit"
	ActionCommandDelay    ActionCommands = "/delay"
	ActionCommandCreate   ActionCommands = "/start"
	ActionCommandSettings ActionCommands = "/settings"
	ActionCommandUndo     ActionCommands = "/undo"
)

var ActionToCommandMap = map[Actions]ActionCommands{
	ActionTake:      ActionCommandTake,
	ActionEdit:      ActionCommandEdit,
	ActionDelay:     ActionCommandDelay,
	ActionCreate:    ActionCommandCreate,
	ActionMySetting: ActionCommandSettings,
	ActionUndo:      ActionCommandUndo,
}

type SendMessageButtons struct {
	Take   bool
	Edit   bool
	Delay  bool
	Create bool
	Undo   bool
}
