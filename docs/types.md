# Types

## Introduction
Ethylene sports a variety of Composite and primitive types.

## Composite Types
### User Defined Types
Enums, Structs, Classes,Abstracts, Interfaces are covered [here](/docs/declarations.md).
They consist of mutable data in fields, as well as callable methods.

Uniquely, Structs, Classes, and Abstracts can be tailed.
Abstracts that are tailed cannot have addition members added when extended.
This is equivalent to the [C struct hack](https://www.geeksforgeeks.org/struct-hack/).
```
import std("io");

class Class[T]~ {
    var T[~] ts;

    public static void(str[]) main = (args) {
        var Class[str]~1 c1;
        c1.ts[0] = "Hello World!";
        io.StdOut.Write(c1.ts[0]);

        var *Class[str]~ cPtr = new(Class[str]~, 1, 1);
        cPtr.ts[0] = "Hello World!";
        io.StdOut.Write(cPtr.ts[0]);
    }
} 
```

### Functions
Functions are callable types used for methods, members, and variables.
They are made up of a typeID and instruction pointer.
```
static var int(int, int) multiply = lambda int(int, int): (x,y) {
    return x * y;
}

public static fun void(str[]) main = (args) {
    var void() f = lambda void() : () {}
    f();
    var int z = multiply(1,2)
}

fun int(int,int) Add = (x, y) {
    return x + y
}
```

### Arrays
### Pointers
### Threads
### Iterators

## Primitive Types
### Integers
### Floats
### Words
### Characters
### Strings
### Booleans
### Void
