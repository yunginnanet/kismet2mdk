package data

import (
	"github.com/bytedance/sonic"
)

func GetStations(d *Device) []string {
	seen := make(map[string]bool)
	st := make([]string, 0, len(d.Dot11.AssociatedClientMap))
	for k, _ := range d.Dot11.AssociatedClientMap {
		if b, ok := seen[k]; ok && b {
			continue
		}
		seen[k] = true
		st = append(st, k)
	}
	return st
}

func Parse(b []byte) (*Device, error) {
	d := NewDevice()
	if err := sonic.Unmarshal(b, d); err != nil {
		return nil, err
	}
	return d, nil
}
