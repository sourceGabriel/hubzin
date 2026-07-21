package firmware

import (
	"time"
)

type ClockWidget struct {
	value string
}

func NewClockWidget() *ClockWidget { return &ClockWidget{} }
func (w *ClockWidget) ID() string  { return "clock" }
func (w *ClockWidget) Update(_ *Device) error {
	w.value = time.Now().Format("15:04:05")
	return nil
}
func (w *ClockWidget) View() map[string]interface{} { return map[string]interface{}{"time": w.value} }

type DateWidget struct {
	value string
}

func NewDateWidget() *DateWidget { return &DateWidget{} }
func (w *DateWidget) ID() string { return "date" }
func (w *DateWidget) Update(_ *Device) error {
	w.value = time.Now().Format("2006-01-02")
	return nil
}
func (w *DateWidget) View() map[string]interface{} { return map[string]interface{}{"date": w.value} }

type RemoteWidget struct {
	id   string
	data interface{}
}

func NewRemoteWidget(id string) *RemoteWidget { return &RemoteWidget{id: id} }
func (w *RemoteWidget) ID() string            { return w.id }
func (w *RemoteWidget) Update(d *Device) error {
	snap, err := d.FetchSnapshot(w.id)
	if err != nil {
		if cached, ok := d.CachedSnapshot(w.id); ok {
			w.data = cached.Data
			return nil
		}
		return err
	}
	w.data = snap.Data
	return nil
}
func (w *RemoteWidget) View() map[string]interface{} { return map[string]interface{}{"data": w.data} }
