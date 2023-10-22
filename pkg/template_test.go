package pkg

import "testing"

func Test_SyncTemplates(t *testing.T) {
	err := SyncTemplates("~/.jumpstart/templates")
	if err != nil {
		t.Error(err)
	}
}
