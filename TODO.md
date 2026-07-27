* thunks should wrap values, not states
* states should contain thunks instead of values
* make calls lazy by wrapping them in an action that returns a state with a thunk
* make lists lazy by wrapping them in a thunk
* prune, sort conversion functions in states
