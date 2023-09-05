# Types
## Kinds of Types

### Primitives
* Void: 0b
  * Symbol: `void`
* Bool: 1b
  * Example Values: `true`, `false`
  * Symbol: `bool`
  * Methods: `str() string`, `Arr[bit8]()uint8s`
  * Operations: `&&`, `&`, `||`, `|`, `!`, `==`, `!=`, `<=>`, `<`, `>`, `<=`, `>=`
* Character: 1b to 4b
  * Example Values: `'a'`, `'/0'`, `'/u1234'`
  * Symbols: `char8`, `char16`, `char32`, `char`
  * Methods: `str() string`, `Arr[bit8]()uint8s`, `bool() isAlpha`, `bool() isDigit`
  * Operations: `+`, `-`, `==`, `!=`, `<=>`, `<`, `>`, `<=`, `>=`
* String: ?b + 8b
  * Example Values: `""`, `"abcde"`, `"abc\n"`
  * Symbols: `str8`, `str16`, `str32`, `str`
  * Methods: `str() string`, `Arr[bit8]()uint8s`
  * Operations: `+`, `[]`, `==`, `!=`, `<=>`, `<`, `>`, `<=`, `>=`
* Integer: 1b to 8b
  * Example Values: `0`, `123`, `0xabcd`
  * Symbols: `int8`, `int16`, `int32`, `int64`, `int`
  * Methods: `str() string`, `Arr[bit8]()uint8s`, `int8() sign`
  * Operations: `+`, `-`, `*`, `/`, `%`, `==`, `!=`, `<=>`, `<`, `>`, `<=`, `>=`
* Floating Point: 1b to 8b
  * Example Values: `0.0`, `1.0e-10`, `3.1415`
  * Symbols: `flt8`, `flt16`, `flt32`, `flt64`, `flt`
  * Methods: `str() string`, `Arr[bit8]()uint8s`, `int8() sign`, `bit16() exponent`, `bit64() mantissa`
  * Operations: `+`, `-`, `*`, `/`, `==`, `!=`, `<=>`, `<`, `>`, `<=`, `>=`
* Fixed Point: 1b to 8b
  * Example Values: `0'0`, `3'14`
  * Symbols: `fix8`, `fix16`, `fix32`, `fix64`, `fix`
  * Methods: `str() string`, `Arr[bit8]()uint8s`, `int64() whole`, `int64() fraction`, `int8() sign`
  * Operations: `+`, `-`, `*`, `/`, `==`, `!=`, `<=>`, `<`, `>`, `<=`, `>=`
* Bit Vector: 1b to 8b
  * Example Values: `0b0`, `0b10`
  * Symbols: `bit8`, `bit16`, `bit32`, `bit64`, `bit`
  * Methods: `str() string`, `Arr[bit8]()uint8s`, `int64() highOne`, `int64() lowOne`, `int64() numOnes`
  * Operations: `&`, `|`, `^`, `!`, `==`, `!=`, `<=>`, `<`, `>`, `<=`, `>=`, `<<`, `>>`
* Type: (8b)type_ptr
  * Example Values: `static(int)`, `static(Obj)`
  * Symbols: `type`
  * Methods: `int64() byteSize`
  * Operations: `@<`, `@>`, `==`, `!=`, `<=>`, `<`, `>`, `<=`, `>=`

### Compound Types
  * Pointers: 2b + 6b
    * Example Value: `null`
    * Example Types: `int**`, `void*`, `Obj*`
    * Extends: `E*` if `(ptr&@ <: E)`
  * Static Array: ?b
    * Example Types: `int[5]`, `bool[4][5]`
    * Methods: ``E(int64) `[]` ``
    * Cannot be assigned.
  * Function: (2b + 6b)op_ptr + (2b + 6b)env_ptr
    * Example Types: `int(int, int)`, `bool()`, `void(int)`
    * Can capture variables
  * Enum: 8b
    * Example Values: `Direction.UP`, `Day.Monday`
    * Example Types: `Day`, `Direction`
    * Methods: `int64() integer`
  * Struct: ?b
    * Example Types: `Complex`, `FatPtr[E]`
    * Cannot extend a parent type.
  * Interface: n/a
    * Example Types: `Number`, `List[E]`
    * Can only extend Interface types.
  * Abstract: n/a
    * Example Types: `File`, `SortedSet[E]`
    * Can only extend Interface types and one Abstract type.
  * Class: (8b)type_ptr + (8b)up_value + ?b
    * Example Types: `Integer`, `LinkList[E]`
    * Can only extend Interface types and one Abstract or Class type.
  * Tailed Class: same as class
    * Example Types: `TailList~[E]`, `TailList~8[E]`
    * Can only extend Interface types and one Abstract or Class type.
  * Tailed Struct: same as struct
    * Example Types: `Node~[E]`, `Node~8[E]`
    * Cannot extend parent type.
    * Garbage collection `mark` method must be specified to run on internal pointers.

Entity(Ent):
* (static(Obj) @< static(Ent)) == true
* (static(ANY_FUN) @< static(Ent)) == true
* (static(ANY_STRUCT) @< static(Ent)) == false
```
class Ent {
  public fun int64(Ent*) cmp = (that) {
    return this{bit64} == that{bit64};
  }

  public fun bool(Ent*) eq = (that) {
    return this{bit64} == that{bit64};
  }
  
  public fun bit64(Ent*) hash = () {
    return this{bit64};
  }
}

class EntEx {
  static public void(Arr[Str]) main = (args) {
    var Complex c = static(Complex).new(1, 1)*;
    var List[Complex] cs = static(List[Complex]).new(c)*;
    if((c& == cs&) == 
      (c&{bit64} == cs&{bit64} || (c& != null && c&.eq(cs)))) {
      print("ALWAYS HAPPENS");
    } else {
      print("NEVER HAPPENS");
    }

    if(c == cs) {
      print("BYTES EQUAL");
    } else {
      print("BYTES UNEQUAL");
    }
  } 
}
```

Object:
* (static(Obj) @< static(Ent)) == true
* (static(ANY_FUN) @< static(Ent)) == false
* (static(ANY_STRUCT) @< static(Ent)) == false
```
class Obj <: [Ent] {
}
```

Struct:
* (static(ANYTHING) @< static(ANY_STRUCT)) ==> SEMANTIC FAILURE
* (static(ANY_STRUCT) @< static(ANYTHING)) ==> SEMANTIC FAILURE
```
struct SomeStruct {
  public fun int64(ExStruct*) cmp = (that) {
    return this{bit64} == that{bit64};
  }

  public fun bool(ExStruct*) eq = (that) {
    return this{bit64} == that{bit64};
  }
  
  public fun bit64(ExStruct*) hash = () {
    return this{bit64};
  }
}

class StructtEx {
  static public void(Arr[Str]) main = (args) {
    var FatPtr[void] fpv = static(FatPtr[void]).new(16, 42{void*});
    var FatPtr[int] fpi = static(FatPtr[int]).new(16, 42{int*});
    if((&fpv == &fpi) == 
      (fpv{bit64} == fpi{bit64} || (fpv != null && fpv.eq(cs)))) { //struct pointer comparison
      print("ALWAYS HAPPENS");
    } else {
      print("NEVER HAPPENS");
    }

    if(fpv == fpi) { //struct value comparison
      print("BYTES EQUAL");
    } else {
      print("BYTES UNEQUAL");
    }
  } 
}
```

Arrays
```
var int[5] is;
var int* is = heap(int, 5);
```

but multidimensional

```
var int[5][5] is1;
var int[25] is2;
//IS NOT
var int** is = heap(int*, 5);
for(var int i = 0; i < 5; i++)
  is[i] = heap(int, 5);
```



Types

typedef struct {
  uint64 funcCount
  uint64[0] funcIntructionPointer
  Map[uint32<<32|uint32, uint64] interfaceFunc#->funcPtr
} ObjType

typedef struct {
  uint32 typeID
  bool lock
  uint64 upValue : 48
  uint64[0] members
} ObjInstance