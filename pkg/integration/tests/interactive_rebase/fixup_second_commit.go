package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FixupSecondCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Fixup the second commit into the first (initial)",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(3)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
			).
			NavigateToListItem(Contains("commit 02")).
			Press(keys.Commits.MarkCommitAsFixup).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Fixup")).
					Content(Equals("Are you sure you want to 'fixup' this commit? It will be merged into the commit below")).
					Confirm()
			}).
			Lines(
				Contains("commit 03"),
				Contains("commit 01").IsSelected(),
			)

		t.Views().Main().
			Content(Contains("commit 01")).
			Content(DoesNotContain("commit 02")).
			Content(Contains("+file01 content")).
			Content(Contains("+file02 content"))
	},
})
