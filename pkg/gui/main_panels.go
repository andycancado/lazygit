package gui

import "os/exec"

type viewUpdateOpts struct {
	title string

	// awkwardly calling this noWrap because of how hard Go makes it to have
	// a boolean option that defaults to true
	noWrap bool

	highlight bool

	task updateTask
}

type refreshMainOpts struct {
	main      *viewUpdateOpts
	secondary *viewUpdateOpts
}

// constants for updateTask's kind field
const (
	RENDER_STRING = iota
	RENDER_STRING_WITHOUT_SCROLL
	RUN_FUNCTION
	RUN_COMMAND
	RUN_PTY
)

type updateTask interface {
	GetKind() int
}

type renderStringTask struct {
	str string
}

func (t *renderStringTask) GetKind() int {
	return RENDER_STRING
}

func (gui *Gui) createRenderStringTask(str string) *renderStringTask {
	return &renderStringTask{str: str}
}

type renderStringWithoutScrollTask struct {
	str string
}

func (t *renderStringWithoutScrollTask) GetKind() int {
	return RENDER_STRING_WITHOUT_SCROLL
}

func (gui *Gui) createRenderStringWithoutScrollTask(str string) *renderStringWithoutScrollTask {
	return &renderStringWithoutScrollTask{str: str}
}

type runCommandTask struct {
	cmd *exec.Cmd
}

func (t *runCommandTask) GetKind() int {
	return RUN_COMMAND
}

func (gui *Gui) createRunCommandTask(cmd *exec.Cmd) *runCommandTask {
	return &runCommandTask{cmd: cmd}
}

type runPtyTask struct {
	cmd *exec.Cmd
}

func (t *runPtyTask) GetKind() int {
	return RUN_PTY
}

func (gui *Gui) createRunPtyTask(cmd *exec.Cmd) *runPtyTask {
	return &runPtyTask{cmd: cmd}
}

type runFunctionTask struct {
	f func(chan struct{}) error
}

func (t *runFunctionTask) GetKind() int {
	return RUN_FUNCTION
}

func (gui *Gui) createRunFunctionTask(f func(chan struct{}) error) *runFunctionTask {
	return &runFunctionTask{f: f}
}

func (gui *Gui) runTaskForView(viewName string, task updateTask) error {
	switch task.GetKind() {
	case RENDER_STRING:
		specificTask := task.(*renderStringTask)
		return gui.newStringTask(viewName, specificTask.str)

	case RENDER_STRING_WITHOUT_SCROLL:
		specificTask := task.(*renderStringWithoutScrollTask)
		return gui.newStringTaskWithoutScroll(viewName, specificTask.str)

	case RUN_FUNCTION:
		specificTask := task.(*runFunctionTask)
		return gui.newTask(viewName, specificTask.f)

	case RUN_COMMAND:
		specificTask := task.(*runCommandTask)
		return gui.newCmdTask(viewName, specificTask.cmd)

	case RUN_PTY:
		specificTask := task.(*runPtyTask)
		return gui.newPtyTask(viewName, specificTask.cmd)
	}

	return nil
}

func (gui *Gui) refreshMain(opts refreshMainOpts) error {
	mainView := gui.getMainView()
	secondaryView := gui.getSecondaryView()

	if opts.main != nil {
		mainView.Title = opts.main.title
		mainView.Wrap = !opts.main.noWrap
		mainView.Highlight = opts.main.highlight // TODO: see what the default should be

		if err := gui.runTaskForView("main", opts.main.task); err != nil {
			gui.Log.Error(err)
			return nil
		}
	}

	gui.splitMainPanel(opts.secondary != nil)

	if opts.secondary != nil {
		secondaryView.Title = opts.secondary.title
		secondaryView.Wrap = !opts.secondary.noWrap
		mainView.Highlight = opts.main.highlight // TODO: see what the default should be
		if err := gui.runTaskForView("secondary", opts.secondary.task); err != nil {
			gui.Log.Error(err)
			return nil
		}
	}

	return nil
}

func (gui *Gui) splitMainPanel(splitMainPanel bool) {
	gui.State.SplitMainPanel = splitMainPanel

	// no need to set view on bottom when splitMainPanel is false: it will have zero size anyway thanks to our view arrangement code.
	if splitMainPanel {
		_, _ = gui.g.SetViewOnTop("secondary")
	}
}

func (gui *Gui) isMainPanelSplit() bool {
	return gui.State.SplitMainPanel
}
