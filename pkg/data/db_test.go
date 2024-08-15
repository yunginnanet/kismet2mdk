package data

import (
	"os"
	"testing"
)

func TestOpenKismetDatabase(t *testing.T) {
	target := os.Getenv("KISMET_TEST_DB")
	if target == "" {
		t.Skip("missing env: 'KISMET_TEST_DB'")
		return
	}
	db, err := OpenKismetDatabase(target)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Run("find related macs", func(t *testing.T) {
		targetMAC := os.Getenv("KISMET_TEST_MAC")
		if targetMAC == "" {
			t.Skip("missing env: 'KISMET_TEST_MAC'")
			return
		}
		related, rerr := db.FindRelatedMacs(targetMAC)
		if rerr != nil {
			t.Fatal(rerr.Error())
		}
		if len(related) == 0 {
			t.Fatal("no related macs found, may be a bad target MAC, please use a MAC with known results")
		}
		seen := make(map[string]bool)
		for _, mac := range related {
			if _, ok := seen[mac]; ok {
				t.Errorf("duplicate MAC: %s", mac)
			}
			seen[mac] = true
			t.Log(mac)
		}
	})
}
