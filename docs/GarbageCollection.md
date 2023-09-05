# Garbage Collection in Ethylene
## Entities, Functions, Objects, Raw Types, Structs, Primitives and You!
Ethylene has many different data types.

### Entities 
  * Are either:
    * *Objects*, which are user defined composite types that are reaped by garbage collection.
    * or *Functions*, which are user defined instructions.
  * If a pointer to them is traceable, they are marked and reaped during garbage collection.

Objects can be *escaped* as nonpointer values, which means they are part of an outer composite type.

```
class Inner {
    int i;
}

class OuterEscape {
    Inner i;
}

class OuterNoEscape {
    *Inner i;
}
```

For example, in the above we wouldn't reap `Inner i` before its containing `OuterEscape`, so it's reap is handled elsewhere and its marking actually marks its containing object.

### Raw Types
  * Are either:
    * *Primitives*, like int, flt, char, bool
    * or *Structs*, which are user defined composite types that are not reaped by garbage collection.
    * If a pointer to them is traceable, the garbage collector won't do squat unless specifically told.

Structs can contain traceable classes.

```
class Markable {
}

struct HasTraceableMember {
    *Markable m; // Must be GC'd!
}
```

As such, HasTraceableMember is checked by the garbage collector as well, but it isn't reaped like collected itself.

```
class SmartClass {
    *int is; // a c-style array, isn't handled ordinarily by GC

    public fun void() mark = () {
        if(is != null)
            delete is; // handle it during GC
    }
}
```

We can also override garbage collection logic through a `mark` method to handle more difficult values.