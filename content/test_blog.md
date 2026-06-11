---
Title: the-test
Tags: test blog
Date: 2026-06-08
Desc: A test blog
---

This is a test blog entry

# This is a header
foo blah

## This is a smaller header?
feet

### Even smaller?

```zig
const std = @import("std");

fn main(init: std.Init) !void {
    const io = init.io;
    std.debug.print("This is a test of highlighting and code embed", .{});
}
```
