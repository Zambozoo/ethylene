# :=o=o=: ETHYLENE :=o=o=:
Check out [h2c2h2.com](http://h2c2h2.com) for out of date resources.


## Compilation Steps
There's a multipass compiler. I'm currently working on 
1. Project File Parsing
    <ol type="a">
    <li>Syntax Parsing</li>
    <li>Semantic Parsing</li>
    </ol>
2. Syntax Parsing:
    <ol type="a">
    <li>CFG Assurance</li>
    <li>Namespace Trie Construction</li>
    </ol>
3. Semantic Parsing
    <ol type="a">
    <li>Parent-Child Assurance</li>
    <li>Cyclical-Avoidance Assurance</li>
    <li>Type Assurance</li>
    <li>Bytecode Generation</li>
    </ol>
4. Bytecode Interpretation

## Main Goal
Cycles and redundant edges are not so fun for a GC.
I mean, they aren't terrible, but they can lead to some redundant GC checks.
My language specification enforces a tracing garbage collector that permits the user to specify 'marking' functionality for their objects if desired.