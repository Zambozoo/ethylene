# Lexical Grammar
## Keywords
```
IMPORT → `import`

CLASS     → `class`
ABSTRACT  → `abstract`
INTERFACE → `interface`
STRUCT    → `struct`
ENUM      → `enum`

PUBLIC    → `public`
PRIVATE   → `private`
PROTECTED → `protected`

STATIC    → `static`
NATIVE    → `native`
VIRTUAL   → `virtual`

FUN       → `fun`
VAR       → `var`

IF       → `if`
ELSE     → `else`
FOR      → `for`
RETURN   → `return`
LABEL    → `label`
BREAK    → `break`
CONTINUE → `continue`
DELETE   → `delete`

ASYNC          → `async`
WAIT           → `wait`
LAMBDA         → `lambda`
TYPE_KEYWORD   → `type`
NEW            → `new`
PANIC          → `panic`
PRINT          → `print`
TRUE           → `true`
FALSE          → `false`
THIS           → `this`
SUPER          → `super`
NULL           → `null`

TYPEINT  → `int`
TYPEFLT  → `flt`
TYPECHAR → `char`
TYPESTR  → `str`
TYPEBOOL → `bool`
TYPEVOID → `void`
TYPEWORD → `word`
```

## Symbols
```
INC → `++`
DEC → `--`

PLUS   → `+`
MINUS  → `-`
STAR   → `*`
DIVIDE → `/`
MODULO → `%`

AND                → `&&`
OR                 → `||`
BITAND             → `&`
BITOR              → `|`
BITXOR             → `^`
BANG               → `!`
SHIFTRIGHT         → `>>`
SHIFTLEFT          → `<<`
SHIFTUNSIGNEDRIGHT → `>>>`

LESSTHAN         → `<`
LESSTHANEQUAL    → `<=`
GREATERTHAN      → `>`
GREATERTHANEQUAL → `>=`
SPACESHIP        → `<=>`
EQUAL            → `==`
BANGEQUAL        → `!=`
ASSIGN           → `=`
NULLARY          → `??`

QUESTIONMARK → `?`
COLON        → `:`
SUBTYPE      → `<:`
SUPERTYPE    → `:>`

PERIOD    → `.`
HASHTAG   → `#`
AT        → `@`
DOLLAR    → `$`
COMMA     → `,`
SEMICOLON → `;`
TILDE     → `~`

LEFTPAREN    → `(`
RIGHTPAREN   → `)`
LEFTBRACKET  → `[`
RIGHTBRACKET → `]`
LEFTBRACE    → `{`
RIGHTBRACE   → `}`
```

## Literals
```
BIN → 0|1
DEC → [0-9]
HEX → [0-9A-F]

INT_DEC    → [1-9] DEC*
INT_HEX    → [1-9A-F] HEX*
INT_BIN    → `1` BIN*
INTEGER    → 0 | INT_DEC | `0x` (INT_HEX|0) | `0b` (INT_BIN|0)

FLT_DEC → (0|INT_DEC) (((`.` INT_DEC)? ([eE] `-`? INT_DEC))|(`.` INT_DEC))
FLT_HEX → (0|INT_HEX) (((`.` INT_HEX)? ([pP] `-`? INT_DEC))|(`.` INT_HEX))
FLOAT   →  FLT_DEX | `0x` FLT_HEX

STRING     →  `"` ([^\]|(`\\` ([ntlr\0"] | `u`HEX)))* `"`
CHARACTER  →  `'` ([^\]*|(`\\` ([ntlr\0'] | `u`HEX))) `"`
IDENTIFIER →  [_A-Za-z][0-9_A-Za-z]*
```