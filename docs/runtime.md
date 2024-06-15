# Runtime
## Data Types
The Ethylene runtime is made up of a linked list of threads, with each thread consisting of:
* The operation stack
* The heap
* The instruction array
* Current Instruction
* Previous scope operation stack address
```go
package runtime

// Shared Runtime contains static data shared by all threads
type SharedRuntime struct {
    // isGarbageCollecting is a flag used while garbage collecting to prevent `new` and `delete` race conditions.
    isGarbageCollecting bool

    // instructions is an array of op code bytes.
    instructions    [?]byte

    // stringLiterals is an array of compile time-defined string bytes.
    stringLiterals       []byte
    
    // staticVariables is an array of static variable bytes.
    staticVariables []byte

    // types is an array of user-defined type definitions for classes, functions, etc.
    types       []byte

    // heap is a pointer to an array of bytes used for dynamic allocations.
    heap            *[]byte
}

type Thread struct {
    // operations is the operations double-stack
    operations              *[]uint64
    // operationsFirstTop points to the top of the first operations stack
    operationsFirstTop           uint64
    // operationsSecondTop points to the top of the second operations stack
    operationsSecondTop    uint64

    // currentInstruction points to the current instruction byte in the shared runtime's instruction array.
    currentInstruction          uint64

    // stackMarkableLinkedListHead points to the markable stack object closest to operationsFirstTop
    stackMarkableLinkedListHead uint64
    // previousScopeTop points to the opening of the current scope in the operationsFirst stack
    currentScopeBottom            uint64

    // threadState contains thread information, such as if it has finished.
    threadState byte

    // result stores the Thread's result.
    result [?]byte
}
```

* The operations stack is actually a kind of "double-stack"([see here]([double](https://www.javatpoint.com/implement-two-stacks-in-an-array#:~:text=in%20the%20array.-,Second%20Approach,-In%20this%20approach) )).