package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Push = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Push a commit to a pre-configured upstream",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")

		shell.CloneIntoRemote("origin")

		shell.SetBranchUpstream("master", "origin/master")

		shell.EmptyCommit("two")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().Content(Contains("↑1 repo → master"))

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.Push)

		assertSuccessfullyPushed(t)
	},
})

func assertSuccessfullyPushed(t *TestDriver) {
	t.Views().Status().Content(Contains("✓ repo → master"))

	t.Views().Remotes().
		Focus().
		Lines(
			Contains("origin"),
		).
		PressEnter()

	t.Views().RemoteBranches().
		IsFocused().
		Lines(
			Contains("master"),
		).
		PressEnter()

	t.Views().SubCommits().
		IsFocused().
		Lines(
			Contains("two"),
			Contains("one"),
		)
}
