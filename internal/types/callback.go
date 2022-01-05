package types

// Callback method signature for structured messages
type Callback func(Message) Message

// StructuredCallback method signature for structured messages
type StructuredCallback func(StructuredMessage) StructuredMessage
