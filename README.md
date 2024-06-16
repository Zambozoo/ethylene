# Ethylene Programming Language
## Introduction
Ethylene is a 
[C-like](https://en.wikipedia.org/wiki/Category:C_programming_language_family), 
[statically typed](https://en.wikipedia.org/wiki/Static_type_system),
[general purpose](https://en.wikipedia.org/wiki/General-purpose_language), 
[imperative](https://en.wikipedia.org/wiki/Imperative_programming), 
[procedural]([procedural](https://en.wikipedia.org/wiki/Procedural_programming)), 
[object-oriented](https://en.wikipedia.org/wiki/Object-oriented_programming), 
[part garbage collected](https://en.wikipedia.org/wiki/Garbage_collection_(computer_science)), 
[part manual memory managed](https://en.wikipedia.org/wiki/Manual_memory_management),
and [compiled](https://en.wikipedia.org/wiki/Compiler)
 [programming language](https://en.wikipedia.org/wiki/Programming_language).

## Getting Started
It's easiest to go look at the [docs](/docs/) and [example](/example/), but here's a quick rundown.

### Folder Structure
Each Ethylene project needs an `eth.yaml` project file, a `pkgs` folder, and some `*.eth` files.
```
root
|-- eth.yaml
|-- pkgs
|   `-- pkg1~0.0.0.zip
`-- src
    `-- Main.eth
```

### Project File
The project file contains project and package details.
The `std` package is reserved for the standard library.
```yaml
name: "example"
version: "0.0.0"
packages:
  std:  "0.0.0"
  pkg1: "0.0.0"
```

### Ethylene Code Files
Ethylene code files are structured similarly to Java.
```
import std("io");

class Complex {
    var flt64 real;
    var flt64 imaginary;

    fun void(flt64, flt64) New = (real, imaginary) {
        this.real = real;
        this.imaginary = imaginary;
    }

    fun Complex(Complex, Complex) Add = (other) {
        var Complex result;
        result.New(this.real + other.real, this.imaginary + other.imaginary);
        return result;
    }

    static fun void([]str) Main = (args) {
        var Complex c1; c1.New(1, 2); // 1+2i
        var Complex c2; c2.New(2, 3); // 2+3i
        var Complex c3 = c1.Add(c2);  // 3+5i
        
        io.StdOut.Write(c3.real{str} + "+" + c3.imaginary{str} + "i");
    }
}