# Ethylene Programming Language
## Description
Ethylene is a work in progress programming language.
It provides a hybrid approach to automated garbage collection and manual memory management.

## Example
Recall that a [Linked List](https://en.wikipedia.org/wiki/Linked_data_structure) consists of nodes chained together like so:
```java
class LinkedList[E] {
  class Node[E] {
    var Node[E]* next;
    var E e;

    // Methods...
  }

  var Node[E]* head;
  var int64 length;

  // Methods...
}
```

A typical [Doubly Linked List](https://en.wikipedia.org/wiki/Doubly_linked_list) contains and extra field for the previous node as well.
```java
class DoublyLinkedList[E] {
  class Node[E] {
    var Node[E]* next;
    var Node[E]* prev;
    var E e;

    // Methods...
  }

  var Node[E]* head;
  var int64 length;

  // Methods...
}
```

An [XOR-Linked List](https://en.wikipedia.org/wiki/XOR_linked_list) takes advantage of the XOR operation's properties to combine both the `next` and `prev` fields into one.
```java
class XORLinkedList[E] {
  class Node[E] {
    var word64 next_xor_prev;
    var E e;

    // Methods...
  }

  var Node[E]* head;
  var int64 length;

  // Methods...
}
```

Unlike the `LinkedList` and `DoublyLinkedList`, a naive garbage collector would be unable to traverse the `XORLinkedList`'s nodes.
This would lead to an early reaping of memory that is currently in use.
Ethylene allows for both the `mark` and `reap` methods to be defined by the user, allowing for the `XORLinkedList` to function properly.

```java
class XORLinkedList[E] {
  // Node definition, other fields...

  // mark marks each element in the list and the list itself.
  public fun void() mark = () {
    var Node[E]* cur = this.head;
    var word64 prev = null{word64};

    for(var int i = 0; i < this.length; i++) {
      var Node[E]* next = (cur.next_xor_prev ^ prev){Node[E]*};
      prev = cur{word64};

      mark(cur.e);
      cur = next;
    }

    mark(this);
  }

  // reap frees each node in the list and the list itself.
  public fun void() reap = () {
    var Node[E]* cur = this.head;
    var word64 prev = null{word64};

    for(var int i = 0; i < this.length; i++) {
      var Node[E]* next = (cur.next_xor_prev ^ prev){Node[E]*};
      prev = cur{word64};

      del(cur);
      cur = next;
    }

    del(this);
  }
}
```

## Special Considerations
As garbage collection can happen at any time, reachable objects must always be prepared for the mark/sweep algorithm whenever heap allocations occur.
