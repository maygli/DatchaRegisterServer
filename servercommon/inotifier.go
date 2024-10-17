package servercommon

type INotifier interface {
	Notify(event []byte)
}
