## Luce

This collection of packages is the result of me (Adam Colton) reusing code. Over
the course of several smaller projects I kept reusing the same code snippets.
I'd find bugs and need to go update many projects. I started putting them all
in this repo and it sort of took on a life of it's own

### Current Status

Because this repo was initially just a place for me to dump shared code, it was
a mess for a long time. Since early 2023 I've been working to clean it up,
which has involved significantly re-writing history. As of early 2024, I'm
almost done with that process.

### Buffers

In most functions that would need to allocate memory a buffer argument is
supplied. It is always safe to use nil for a buffer. This pattern can avoid
generating garbage often enough that it is useful and (by simply supplying nil)
it is easy to ignore.