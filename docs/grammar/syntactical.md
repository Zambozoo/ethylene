# Syntactical Grammar
## Declarations
```
GENERIC_ARG → IDENTIFIER (SUBTYPE|SUPERTYPE) (IDENTFIER | LEFTBRACKET TYPE (COMMA TYPE)* RIGHTBRACKET)
GENERIC_DECL → IDENTIFIER (LEFTBRACKET GENERIC_ARG (COMMA GENERIC_ARG)* RIGHTBRACKET)?

PARENTS → SUBTYPE (TYPE | LEFTBRACKET TYPE (COMMA TYPE)* RIGHTBRACKET)

DECL_CLASS     → CLASS GENERIC_DECL TILDE? PARENTS? LEFTBRACE NONVIRTUAL_FIELD* RIGHTBRACE
DECL_ABSTRACT  → ABSTRACT GENERIC_DECL TILDE? PARENTS? LEFTBRACE FIELD* RIGHTBRACE
DECL_INTERFACE → INTERFACE GENERIC_DECL PARENTS? LEFTBRACE INTERFACE_FIELD* RIGHTBRACE
DECL_STRUCT    → STRUCT GENERIC_DECL TILDE? LEFTBRACE NON_VIRTUAL_FIELD* RIGHTBRACE
DECL_ENUM      → CLASS IDENTIFIER LEFTBRACE (ENUM_FIELD (COMMA ENUM_FIELD)* SEMICOLON)? NONVIRTUAL_FIELD* RIGHTBRACE
```

## Fields
```
ACCESS_MODIFIER  → PUBLIC|PRIVATE|PROTECTED

NONVIRTUAL_FIELD → NONVIRTUAL_METHOD | FIELD_MEMBER | FIELD_DECL
INTERFACE_FIELD  → STATIC_MEMBER | VIRTUAL_METHOD | STATIC_METHOD | FIELD_DECL
FIELD            → METHOD | MEMBER | FIELD_DECL

NOBODY_METHOD     → FUNC TYPE IDENTIFIER SEMICOLON
VIRTUAL_METHOD    → {ACCESS_MODIFIER? VIRTUAL} NOBODY_METHOD
NATIVE_METHOD     → {ACCESS_MODIFIER? {STATIC? NATIVE}} NOBODY_METHOD
BODY_METHOD       → FUNC TYPE IDENTIFIER ASSIGN LEFTPAREN (IDENTIFIER (COMMA IDENTIFIER)*)? STMT_BLOCK
STATIC_METHOD     → ({ACCESS_MODIFIER? {STATIC NATIVE}} NOBODY_METHOD)| {ACCESS_MODIFIER? STATIC} BODY_METHOD
NONVIRTUAL_METHOD → NATIVE_METHOD| {ACCESS_MODIFIER? STATIC?} BODY_METHOD
METHOD            → NONVIRTUAL_METHOD|VIRTUAL_METHOD

STATIC_MEMBER    → {ACCESS_MODIFIER? STATIC} FIELD_MEMBER
NONSTATIC_MEMBER → ACCESS_MODIFIER? VAR TYPE IDENTIFIER (ASSIGN EXPR)? SEMICOLON
MEMBER           → STATIC_MEMBER|NONSTATIC_MEMBER

FIELD_DECL   → ACCESS_MODIFIER? (DECL_CLASS|DECL_ABSTRACT|DECL_INTERFACE|DECL_STRUCT|DECL_ENUM)
FIELD_ENUM   → IDENTIFIER (PERIOD EXPR_CALL)?
```

## Statements
```
STMT_BLOCK    → LEFTBRACE STMT* RIGHTBRACE
STMT_BREAK    → BREAK IDENTIFIER? SEMICOLON
STMT_CONTINUE → CONTINUE IDENTIFIER? SEMICOLON
STMT_DELETE   → DELETE LEFTPAREN EXPR RIGHTPAREN SEMICOLON
STMT_EXPR     → EXPR SEMICOLON
STMT_FOR0     → FOR STMT
STMT_FOR1     → FOR LEFTPAREN EXPR RIGHTPAREN STMT (ELSE STMT)?
STMT_FOR3     → FOR LEFTPAREN (VAR|EXPR SEMICOLON| SEMICOLON) EXPR? SEMICOLON EXPR? RIGHTPAREN STMT (ELSE STMT)?
STMT_FOREACH  → FOR LEFTPAREN VAR TYPE IDENTIFIER COLON EXPR RIGHTPAREN STMT (ELSE STMT)?
STMT_IF       → IF LEFTPAREN EXPR RIGHTPAREN STMT (ELSE STMT)?
STMT_LABEL    → LABEL IDENTIFIER COLON STMT
STMT_PANIC    → PANIC LEFTPARENT EXPR RIGHTPAREN SEMICOLON
STMT_PRINT    → PRINT LEFTPARENT EXPR RIGHTPAREN SEMICOLON
STMT_RETURN   → RETURN EXPR? SEMICOLON
STMT_VAR      → VAR TYPE IDENTIFIER (ASSIGN EXPR)? SEMICOLON
```

## Expressions
```
EXPR_ACCESS     → EXPR LEFTBRACKET EXPR RIGHTBRACKET
EXPR_ASSIGN     → EXPR_TERNARY (EQUAL EXPR_TERNARY)?
EXPR_ASYNC      → ASYNC EXPR
EXPR_BITAND     → EXPR_EQUAL (BITAND EXPR_EQUAL)? 
EXPR_BITOR      → EXPR_XOR (BITOR EXPR_XOR)?
EXPR_BITXOR     → EXPR_BITAND (BITOR EXPR_BITAND)?
EXPR_CALL       → EXPR LEFTPAREN (EXPR (COMMA EXPR))? RIGHTPAREN
EXPR_CAST       → EXPR LEFTBRACE (EXPR (COMMA EXPR))? RIGHTBRACE
EXPR_COMPARE    → EXPR_SHIFT ((GREATER|GREATERTHAN|LESS|LESSTHAN|COMPARE|SUBTYPE|SUPERTYPE) EXPR_SHIFT)? 
EXPR_EQUAL      → EXPR_COMPARE ((EQUAL|NOTEQUAL) EXPR_COMPARE)? 
EXPR_FACTOR     → EXP_UNARYPOST ((STAR|DIVIDE|MODULO) EXPR_UNARYPOST)? 
EXPR_FIELD      → IDENTIFIER (PERIOD IDENTIFIER)*
EXPR_TYPE       → TYPE_KEYWORD LEFTPAREN TYPE RIGHTPAREN
EXPR_TYPEFIELD  → TYPE_KEYWORD LEFTPAREN COMPOSITE_TYPE RIGHTPAREN
EXPR_UNARYPRE   → (INCREMENT|DECREMENT|STAR|MINUS|BANG|AT|HASH) EXPR_PRIMARY
EXPR_UNARYPOST  → EXPR_UNARYPRE (INCREMENT|DECREMENT)?
EXPR_IDENTIFIER → IDENTIFIER
EXPR_LAMBDA     → LAMBDA TYPE COLON LEFTPARENT (IDENTIFIER (COMMA IDENTIFIER)*)? STMT
EXPR_LITERAL    → (INTEGER|FLOAT|STRING|CHAR|TRUE|FALSE|NULL|THIS|SUPER)
EXPR_AND        → EXPR_BITOR (NULLARY EXPR_BITOR)?
EXPR_OR         → EXPR_AND (NULLARY EXPR_AND)?
EXPR_NEW        → NEW LEFTPAREN TYPE (COMMA EXPR (COMMA EXPR)?)? RIGHTPAREN 
EXPR_NULLARY    → EXPR_OR (NULLARY EXPR_OR)?
EXPR_PRIMARY    → EXPR_WAIT|EXPR_LAMBDA|EXPR_IDENTIFIER|EXPR_FIELD|EXPR_CALL|EXPR_ACCESS|EXPR_TYPE|EXPR_TYPEFIELD|EXPR_NEW|EXPR_CAST|EXPR_LITERAL
EXPR_SHIFT      → EXPR_TERM ((RIGHTSHIFT|LEFTSHIFT|SHIFTUNSIGNEDRIGHT) EXPR_TERM)? 
EXPR_TERM       → EXPR_FACTOR ((PLUS|MINUS) EXPR_FACTOR)? 
EXPR_TERNARY    → EXPR_NULLARY (QUESTIONMARK EXPR_NULLARY COLON EXPR_NULLARY)?
EXPR_WAIT       → WAIT LEFTPAREN EXPR RIGHTPAREN
```

## Types
```
TYPE_ARRAY     → TYPE LEFTBRACKET RIGHTBRACKET
TYPE_COMPOSITE → IDENTIFIER (PERIOD IDENTIFIER)*
TYPE_FUNCTION  → TYPE LEFTPAREN (TYPE (COMMA TYPE)*)? RIGHTPAREN
TYPE_GENERIC   → TYPE LEFTBRACKET (TYPE (COMMA TYPE)*)? RIGHTBRACKET
TYPE_POINTER   → TYPE STAR
TYPE_PRIMITIVE → TYPEINT|TYPEFLT|TYPEWORD|TYPECHAR|TYPESTR|TYPEBOOL|TYPEVOID
TYPE_TAILED    → (TYPE_COMPOSITE|TYPE_GENERIC) TILDE INTEGER?
````