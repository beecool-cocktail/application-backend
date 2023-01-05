package viewmodels

type CommandUndoRequest struct {
	// ID of a command
	//
	// In: path
	ID string `uri:"id" binding:"required"`
}
