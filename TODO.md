* thunks should wrap values, not states
* states should contain thunks instead of values
* make calls lazy by wrapping them in an action that returns a state with a thunk
* make lists lazy by wrapping them in a thunk
* FIXME: some Replace calls need to be replaced, they are passing variable bindings and shouldn't
* prune, sort conversion functions in states
