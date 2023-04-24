package magic

type EventHandler func(ev string, data any)
type EventSender func(ev string, data any)

const MountEvent = "magic:mount"
