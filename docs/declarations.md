# Declarations
## Introduction
There are five kinds of declarations in Ethylene:
* Classes
* Abstracts
* Interfaces
* Structs
* Enums

## Classes
Classes are garbage collected, user-defined, composite types.
They can extend a single class or abtract, as well as implement many interfaces.
They cannot have 'virtual' methods.
They can use generic type parameters.

## Abstracts
Abstracts are used for shared functionality between classes.
They can extend a single abtract, as well as implement many interfaces.
They can have 'virtual' and 'concrete' methods.
They can use generic type parameters.
They only exist as pointers to classes.

## Interfaces
Interfaces are used for ad hoc polymophism.
They can implement many interfaces.
They cannot have non-static 'concrete' methods.
They can use generic type parameters.
They only exist as pointers to classes.

## Structs
Structs are user-defined and composite types whose lifecycle is user-managed.
They cannot extend or implement other declarations.
They cannot have 'virtual' methods.
They can use generic type parameters.

## Enums
Enums are static, user-defined, ordinal types.
They cannot be created with `new` at runtime.
They cannot have 'virtual' methods.
They cannot use generic type parameters.