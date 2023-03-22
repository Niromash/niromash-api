package model

import (
	"testing"
)

func TestHasPermission(t *testing.T) {
	user := User{Permissions: []string{"messages.create"}}
	if !user.HasPermission("messages.create") {
		t.Error("User should have permission messages.create")
	}

	if user.HasPermission("messages.delete") {
		t.Error("User should have permission messages.delete")
	}

	if user.HasPermission("messages.*") {
		t.Error("User should have permission messages.*")
	}

	user.Permissions = []string{"messages.*"}
	if !user.HasPermission("messages.create") {
		t.Error("User should have permission messages.create")
	}

	if !user.HasPermission("messages.delete") {
		t.Error("User should have permission messages.delete")
	}

	if !user.HasPermission("messages.*") {
		t.Error("User should have permission messages.*")
	}

	user.Permissions = []string{"*"}
	if !user.HasPermission("messages.create") {
		t.Error("User should have permission messages.create")
	}

	if !user.HasPermission("messages.delete") {
		t.Error("User should have permission messages.delete")
	}

	if !user.HasPermission("messages.*") {
		t.Error("User should have permission messages.*")
	}

	user.Permissions = []string{"messages.create", "messages.delete"}
	if !user.HasPermission("messages.create") {
		t.Error("User should have permission messages.create")
	}

	if !user.HasPermission("messages.delete") {
		t.Error("User should have permission messages.delete")
	}

	if user.HasPermission("messages.*") {
		t.Error("User should have permission messages.*")
	}
}
