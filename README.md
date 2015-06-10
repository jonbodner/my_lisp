# my_lisp
A simple Lisp-like language that I'm implementing in Go.

Writing a Lisp is something I'd always wanted to do and I wanted to use Go in a 
project that was slightly bigger than a programming puzzle.

So far the following functionality is there:
- CONS
- CAR
- CDR
- QUOTE
- LAMBDA
- SETQ
- ATOM
- EQ
- COND
- + - * /
- Infinite precision math (Integers and ratios)

It's a LISP-1 (single namespace for both values and functions). The scoping is static.

Other features that I intend to add (in likely order):
- PROGN
- LET
- Tail Call optimization
- Read/Write environment
- Macros
- CSP functionality (Goroutines, Channels, Select)
- Maps, Sets

As a stretch goal, I'd like to add the ability to generate compiled code as well.

At some point, I'm going to start renaming the special forms and switching to lower case (LAMBDA -> fn, for example).

